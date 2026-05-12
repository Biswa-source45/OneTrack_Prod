package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, you should validate the origin
		return true
	},
}

// WebSocketHandler handles WebSocket connection requests
func WebSocketHandler(hub *Hub, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user information from context (set by auth middleware)
		var userID uint
		var userType string

		// Check if it's a staff user (manager/engineer)
		if userVal, exists := c.Get("user"); exists {
			if user, ok := userVal.(models.User); ok {
				userID = user.ID
				userType = "user"
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
		} else if contactIDVal, exists := c.Get("contact_id"); exists {
			// Customer contact
			if contactID, ok := contactIDVal.(uint); ok {
				userID = contactID
				userType = "contact"
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact in context"})
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("❌ WebSocket upgrade failed: %v", err)
			return
		}

		// Create client instance
		client := &Client{
			ID:           uuid.New().String(),
			UserID:       userID,
			UserType:     userType,
			Conn:         conn,
			Hub:          hub,
			Send:         make(chan []byte, 256),
			lastPingTime: time.Now(),
		}

		// Register client with hub
		hub.register <- client

		// Send initial connection success message
		welcomeMsg := Message{
			Type: "connected",
			Data: map[string]interface{}{
				"client_id": client.ID,
				"user_id":   userID,
				"user_type": userType,
				"message":   "WebSocket connection established",
			},
			Timestamp: time.Now(),
		}

		if welcomeBytes, err := welcomeMsg.ToJSON(); err == nil {
			client.Send <- welcomeBytes
		}

		// Fetch and send current unread count
		go func() {
			var count int64
			db.Model(&models.Notification{}).
				Where("recipient_id = ? AND recipient_type = ? AND is_read = ?", userID, userType, false).
				Count(&count)

			countMsg := Message{
				Type: "count.update",
				Data: map[string]interface{}{
					"unread_count": count,
				},
				Timestamp: time.Now(),
			}

			if countBytes, err := countMsg.ToJSON(); err == nil {
				client.Send <- countBytes
			}
		}()

		log.Printf("🔌 WebSocket connection established: ClientID=%s, UserType=%s, UserID=%d",
			client.ID, userType, userID)

		// Start client goroutines
		go client.WritePump()
		go client.ReadPump()
	}
}

// ToJSON converts Message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"type":"%s","data":%s,"timestamp":"%s"}`,
		m.Type,
		toJSONString(m.Data),
		m.Timestamp.Format(time.RFC3339),
	)), nil
}

// Helper function to convert data to JSON string
func toJSONString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case map[string]interface{}:
		result := "{"
		first := true
		for key, val := range v {
			if !first {
				result += ","
			}
			result += fmt.Sprintf(`"%s":%s`, key, toJSONString(val))
			first = false
		}
		result += "}"
		return result
	case int, int64, uint, uint64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf(`"%v"`, v)
	}
}
