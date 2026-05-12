package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	ID           string          // Unique client identifier
	UserID       uint            // User or Contact ID
	UserType     string          // "user" or "contact"
	Conn         *websocket.Conn // WebSocket connection
	Hub          *Hub            // Reference to hub
	Send         chan []byte     // Buffered channel for outbound messages
	mu           sync.Mutex      // Mutex for thread-safe operations
	lastPingTime time.Time       // Last ping timestamp
}

// Hub maintains active WebSocket connections and broadcasts messages
type Hub struct {
	// Registered clients mapped by UserType:UserID (e.g., "user:5", "contact:10")
	clients map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to specific user
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe map operations
	mu sync.RWMutex
}

// BroadcastMessage represents a message to be sent to specific user(s)
type BroadcastMessage struct {
	UserID      uint   // Target user ID
	UserType    string // Target user type ("user" or "contact")
	MessageType string // Type of message
	Data        interface{}
}

// Message represents a WebSocket message structure
type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256), // Buffered channel
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	log.Println("🚀 WebSocket Hub started")
	
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToUser(message)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := h.getUserKey(client.UserType, client.UserID)
	
	if h.clients[key] == nil {
		h.clients[key] = make(map[*Client]bool)
	}
	
	h.clients[key][client] = true
	
	log.Printf("✅ Client registered: %s (UserType: %s, UserID: %d) - Total connections for this user: %d",
		client.ID, client.UserType, client.UserID, len(h.clients[key]))
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := h.getUserKey(client.UserType, client.UserID)
	
	if clients, ok := h.clients[key]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.Send)
			
			// Clean up empty user entry
			if len(clients) == 0 {
				delete(h.clients, key)
			}
			
			log.Printf("❌ Client unregistered: %s (UserType: %s, UserID: %d) - Remaining connections: %d",
				client.ID, client.UserType, client.UserID, len(clients))
		}
	}
}

// broadcastToUser sends a message to all connections of a specific user
func (h *Hub) broadcastToUser(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	key := h.getUserKey(msg.UserType, msg.UserID)
	clients, ok := h.clients[key]
	
	if !ok || len(clients) == 0 {
		log.Printf("⚠️  No active connections for %s:%d", msg.UserType, msg.UserID)
		return
	}

	// Create message payload
	message := Message{
		Type:      msg.MessageType,
		Data:      msg.Data,
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Error marshaling message: %v", err)
		return
	}

	// Send to all client connections for this user
	successCount := 0
	for client := range clients {
		select {
		case client.Send <- messageBytes:
			successCount++
		default:
			// Client's send buffer is full, skip
			log.Printf("⚠️  Client %s send buffer full, message dropped", client.ID)
		}
	}

	log.Printf("📤 Broadcast to %s:%d - Sent to %d/%d connections (Type: %s)",
		msg.UserType, msg.UserID, successCount, len(clients), msg.MessageType)
}

// BroadcastToUser sends a message to a specific user
func (h *Hub) BroadcastToUser(userID uint, userType, messageType string, data interface{}) {
	h.broadcast <- &BroadcastMessage{
		UserID:      userID,
		UserType:    userType,
		MessageType: messageType,
		Data:        data,
	}
}

// GetActiveConnectionCount returns the number of active connections for a user
func (h *Hub) GetActiveConnectionCount(userID uint, userType string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	key := h.getUserKey(userType, userID)
	if clients, ok := h.clients[key]; ok {
		return len(clients)
	}
	return 0
}

// GetTotalConnectionCount returns total number of active connections
func (h *Hub) GetTotalConnectionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	return total
}

// getUserKey creates a unique key for user identification
func (h *Hub) getUserKey(userType string, userID uint) string {
	return userType + ":" + string(rune(userID))
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// Configure connection
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.mu.Lock()
		c.lastPingTime = time.Now()
		c.mu.Unlock()
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket error for client %s: %v", c.ID, err)
			}
			break
		}

		// Handle incoming messages (ping, pong, etc.)
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("⚠️  Invalid message from client %s: %v", c.ID, err)
			continue
		}

		// Handle different message types
		switch msg.Type {
		case "ping":
			// Respond with pong
			pongMsg := Message{
				Type:      "pong",
				Data:      nil,
				Timestamp: time.Now(),
			}
			if pongBytes, err := json.Marshal(pongMsg); err == nil {
				c.Send <- pongBytes
			}
		default:
			log.Printf("📨 Received message from client %s: Type=%s", c.ID, msg.Type)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second) // Ping interval
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			
			if !ok {
				// Hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
