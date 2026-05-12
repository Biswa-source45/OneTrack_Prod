package routes

import (
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/handlers"
	"github.com/Chinmay-Globx/ticketing-backend/internal/middleware"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	ws "github.com/Chinmay-Globx/ticketing-backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, hub *ws.Hub) *gin.Engine {
	r := gin.Default()
	// CORS: allow Vite dev server and Authorization header
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inject WebSocket hub into notification service globally
	// This ensures all notification service instances can broadcast
	services.GlobalWebSocketHub = hub

	// Apply audit middleware globally to capture request context
	r.Use(middleware.AuditMiddleware())

	r.POST("/refresh", handlers.RefreshToken(db))

	// WebSocket endpoint for real-time notifications
	r.GET("/ws/notifications", handlers.AuthMiddleware(db), ws.WebSocketHandler(hub, db))

	// Auth routes (login endpoints do NOT require JWT)
	r.POST("/login/user", handlers.UserLogin(db))
	r.POST("/login/contact", handlers.ContactLogin(db))

	// All other routes require JWT
	r.POST("/logout", handlers.AuthMiddleware(db), handlers.Logout(db))
	r.POST("/reset-password/user", handlers.AuthMiddleware(db), handlers.ResetUserPassword(db))
	r.POST("/reset-password/contact", handlers.AuthMiddleware(db), handlers.ResetContactPassword(db))

	r.GET("/protected", handlers.AuthMiddleware(db), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "You are authenticated!"})
	})

	r.POST("/accounts", handlers.AuthMiddleware(db), handlers.CreateAccount(db))
	r.GET("/accounts", handlers.AuthMiddleware(db), handlers.GetAccounts(db))
	r.GET("/accounts/:id", handlers.AuthMiddleware(db), handlers.GetAccount(db))
	r.PUT("/accounts/:id", handlers.AuthMiddleware(db), handlers.UpdateAccount(db))
	r.DELETE("/accounts/:id", handlers.AuthMiddleware(db), handlers.DeleteAccount(db))

	r.POST("/contacts", handlers.AuthMiddleware(db), handlers.CreateContact(db))
	r.GET("/contacts", handlers.AuthMiddleware(db), handlers.GetContacts(db))
	r.GET("/contacts/:id", handlers.AuthMiddleware(db), handlers.GetContact(db))
	r.PUT("/contacts/:id", handlers.AuthMiddleware(db), handlers.UpdateContact(db))
	r.DELETE("/contacts/:id", handlers.AuthMiddleware(db), handlers.DeleteContact(db))

	r.POST("/users", handlers.AuthMiddleware(db), handlers.CreateUser(db))
	r.GET("/users", handlers.AuthMiddleware(db), handlers.GetUsers(db))
	r.GET("/users/:id", handlers.AuthMiddleware(db), handlers.GetUser(db))
	r.PUT("/users/:id", handlers.AuthMiddleware(db), handlers.UpdateUser(db))
	r.DELETE("/users/:id", handlers.AuthMiddleware(db), handlers.DeleteUser(db))

	r.POST("/products", handlers.AuthMiddleware(db), handlers.CreateProduct(db))
	r.GET("/products", handlers.AuthMiddleware(db), handlers.GetProducts(db))
	r.GET("/products/:id", handlers.AuthMiddleware(db), handlers.GetProduct(db))
	r.PUT("/products/:id", handlers.AuthMiddleware(db), handlers.UpdateProduct(db))
	r.DELETE("/products/:id", handlers.AuthMiddleware(db), handlers.DeleteProduct(db))

	// Issues
	r.POST("/issues", handlers.AuthMiddleware(db), handlers.CreateIssue(db))
	r.GET("/issues", handlers.AuthMiddleware(db), handlers.GetIssues(db))
	r.PUT("/issues/:id", handlers.AuthMiddleware(db), handlers.UpdateIssue(db))
	r.DELETE("/issues/:id", handlers.AuthMiddleware(db), handlers.DeleteIssue(db))

	// Product Issue CRUD
	r.POST("/product-issues", handlers.AuthMiddleware(db), handlers.CreateProductIssue(db))
	r.GET("/product-issues", handlers.AuthMiddleware(db), handlers.GetProductIssues(db))
	r.GET("/product-issues/:id", handlers.AuthMiddleware(db), handlers.GetProductIssue(db))
	r.PUT("/product-issues/:id", handlers.AuthMiddleware(db), handlers.UpdateProductIssue(db))
	r.DELETE("/product-issues/:id", handlers.AuthMiddleware(db), handlers.DeleteProductIssue(db))

	r.POST("/roles", handlers.AuthMiddleware(db), handlers.CreateRole(db))
	r.GET("/roles", handlers.AuthMiddleware(db), handlers.GetRoles(db))
	r.PUT("/roles/:id", handlers.AuthMiddleware(db), handlers.UpdateRole(db))
	r.DELETE("/roles/:id", handlers.AuthMiddleware(db), handlers.DeleteRole(db))

	r.POST("/designations/users", handlers.AuthMiddleware(db), handlers.CreateUserDesignation(db))
	r.GET("/designations/users", handlers.AuthMiddleware(db), handlers.GetUserDesignations(db))
	r.PUT("/designations/users/:id", handlers.AuthMiddleware(db), handlers.UpdateUserDesignation(db))
	r.DELETE("/designations/users/:id", handlers.AuthMiddleware(db), handlers.DeleteUserDesignation(db))

	r.POST("/designations/contacts", handlers.AuthMiddleware(db), handlers.CreateContactDesignation(db))
	r.GET("/designations/contacts", handlers.AuthMiddleware(db), handlers.GetContactDesignations(db))
	r.PUT("/designations/contacts/:id", handlers.AuthMiddleware(db), handlers.UpdateContactDesignation(db))
	r.DELETE("/designations/contacts/:id", handlers.AuthMiddleware(db), handlers.DeleteContactDesignation(db))

	// Ticket routes
	r.POST("/tickets", handlers.AuthMiddleware(db), handlers.CreateTicketHandler(db))
	r.GET("/tickets", handlers.AuthMiddleware(db), handlers.GetTicketsByContactHandler(db))
	r.GET("/tickets/:id", handlers.GetTicketDetailHandler(db))                                 // Public endpoint for ticket details
	r.GET("/tickets/:id/full", handlers.AuthMiddleware(db), handlers.GetTicketFullDetails(db)) // Enhanced ticket details

	// Ticket attachment routes
	r.POST("/tickets/:id/attachments", handlers.AuthMiddleware(db), handlers.UploadTicketAttachments(db))
	r.GET("/tickets/:id/attachments", handlers.AuthMiddleware(db), handlers.GetTicketAttachments(db))
	r.GET("/attachments/:id/download", handlers.AuthMiddleware(db), handlers.DownloadAttachment(db))
	r.DELETE("/attachments/:id", handlers.AuthMiddleware(db), handlers.DeleteAttachment(db))

	// Ticket Comments routes
	r.POST("/tickets/:id/comments", handlers.AuthMiddleware(db), handlers.CreateTicketComment(db))
	r.GET("/tickets/:id/comments", handlers.AuthMiddleware(db), handlers.GetTicketComments(db))
	r.PUT("/tickets/:id/comments/:comment_id", handlers.AuthMiddleware(db), handlers.UpdateTicketComment(db))
	r.DELETE("/tickets/:id/comments/:comment_id", handlers.AuthMiddleware(db), handlers.DeleteTicketComment(db))

	// Ticket Calls routes
	r.POST("/tickets/:id/calls", handlers.AuthMiddleware(db), handlers.CreateTicketCall(db))
	r.GET("/tickets/:id/calls", handlers.AuthMiddleware(db), handlers.GetTicketCalls(db))
	r.PUT("/tickets/:id/calls/:call_id", handlers.AuthMiddleware(db), handlers.UpdateTicketCall(db))
	r.PATCH("/tickets/:id/calls/:call_id/complete", handlers.AuthMiddleware(db), handlers.CompleteTicketCall(db))
	r.PATCH("/tickets/:id/calls/:call_id/cancel", handlers.AuthMiddleware(db), handlers.CancelTicketCall(db))
	r.PATCH("/tickets/:id/calls/:call_id/close", handlers.AuthMiddleware(db), handlers.CloseTicketCall(db))

	// Ticket Call Attachments routes
	r.POST("/tickets/:id/calls/:call_id/attachments", handlers.AuthMiddleware(db), handlers.UploadCallAttachments(db))
	r.GET("/tickets/:id/calls/:call_id/attachments", handlers.AuthMiddleware(db), handlers.GetCallAttachments(db))
	r.GET("/call-attachments/:attachment_id/download", handlers.AuthMiddleware(db), handlers.DownloadCallAttachment(db))
	r.DELETE("/call-attachments/:attachment_id", handlers.AuthMiddleware(db), handlers.DeleteCallAttachment(db))

	// Ticket Activities routes
	r.GET("/tickets/:id/activities", handlers.AuthMiddleware(db), handlers.GetTicketActivities(db))
	r.GET("/tickets/:id/timeline", handlers.AuthMiddleware(db), handlers.GetTicketTimeline(db))

	// Ticket Approvals routes
	r.POST("/tickets/:id/approvals", handlers.AuthMiddleware(db), handlers.CreateApprovalRequestHandler(db))
	r.GET("/tickets/:id/approvals", handlers.AuthMiddleware(db), handlers.ListApprovalRequestsHandler(db))
	r.PATCH("/tickets/:id/approvals/:approvalId/approve", handlers.AuthMiddleware(db), handlers.ApproveRequestHandler(db))
	r.PATCH("/tickets/:id/approvals/:approvalId/reject", handlers.AuthMiddleware(db), handlers.RejectRequestHandler(db))

	// Manager module routes
	r.GET("/manager/dashboard/stats", handlers.AuthMiddleware(db), handlers.ManagerDashboardStatsHandler(db))
	r.POST("/manager/tickets", handlers.AuthMiddleware(db), handlers.ManagerCreateTicketHandler(db))
	r.GET("/manager/tickets", handlers.AuthMiddleware(db), handlers.ManagerListTicketsHandler(db))
	r.PUT("/manager/tickets/:id", handlers.AuthMiddleware(db), handlers.ManagerEditTicketHandler(db))
	r.DELETE("/manager/tickets/:id", handlers.AuthMiddleware(db), handlers.ManagerDeleteTicketHandler(db))
	r.PATCH("/manager/tickets/:id/status", handlers.AuthMiddleware(db), handlers.ManagerChangeStatusHandler(db))
	r.PATCH("/manager/tickets/:id/assign", handlers.AuthMiddleware(db), handlers.ManagerAssignTicketHandler(db))
	r.GET("/manager/engineers", handlers.AuthMiddleware(db), handlers.ManagerListEngineersHandler(db))
	r.GET("/manager/engineers-with-tickets", handlers.AuthMiddleware(db), handlers.ManagerListEngineersWithTicketCountHandler(db))
	r.GET("/manager/engineers/:id/tickets", handlers.AuthMiddleware(db), handlers.ManagerEngineerTicketsHandler(db))

	// Manager task routes
	r.POST("/manager/tasks", handlers.AuthMiddleware(db), handlers.ManagerCreateTaskHandler(db))
	r.GET("/manager/tasks", handlers.AuthMiddleware(db), handlers.ManagerGetTasksHandler(db))
	r.GET("/manager/tasks/:id", handlers.AuthMiddleware(db), handlers.ManagerGetTaskHandler(db))
	r.PUT("/manager/tasks/:id", handlers.AuthMiddleware(db), handlers.ManagerEditTaskHandler(db))
	r.DELETE("/manager/tasks/:id", handlers.AuthMiddleware(db), handlers.ManagerDeleteTaskHandler(db))

	// Dumped Queries routes (Manager/Admin)
	dumpHandler := handlers.NewDumpedQueryHandler(db)
	r.GET("/manager/dumped-queries", handlers.AuthMiddleware(db), dumpHandler.GetDumpedQueries)
	r.GET("/manager/dumped-queries/:id", handlers.AuthMiddleware(db), dumpHandler.GetDumpedQuery)
	r.DELETE("/manager/dumped-queries/:id", handlers.AuthMiddleware(db), dumpHandler.DeleteDumpedQuery)
	r.PATCH("/manager/dumped-queries/:id/status", handlers.AuthMiddleware(db), dumpHandler.UpdateDumpedQueryStatus)

	// Task comments routes
	r.POST("/tasks/:id/comments", handlers.AuthMiddleware(db), handlers.CreateTaskComment(db))
	r.GET("/tasks/:id/comments", handlers.AuthMiddleware(db), handlers.GetTaskComments(db))
	r.PUT("/tasks/:id/comments/:commentId", handlers.AuthMiddleware(db), handlers.UpdateTaskComment(db))
	r.DELETE("/tasks/:id/comments/:commentId", handlers.AuthMiddleware(db), handlers.DeleteTaskComment(db))

	// Task activities routes
	r.GET("/tasks/:id/activities", handlers.AuthMiddleware(db), handlers.GetTaskActivities(db))

	// Engineer module routes
	r.POST("/engineer/tickets", handlers.AuthMiddleware(db), handlers.EngineerCreateTicketHandler(db))
	r.GET("/engineer/tickets", handlers.AuthMiddleware(db), handlers.EngineerListTicketsHandler(db))
	r.PATCH("/engineer/tickets/:id/status", handlers.AuthMiddleware(db), handlers.EngineerChangeStatusHandler(db))

	// Engineer task routes
	r.GET("/engineer/tasks", handlers.AuthMiddleware(db), handlers.EngineerListTasksHandler(db))
	r.GET("/engineer/tasks/:id", handlers.AuthMiddleware(db), handlers.EngineerGetTaskHandler(db))
	r.PATCH("/engineer/tasks/:id/status", handlers.AuthMiddleware(db), handlers.EngineerChangeTaskStatusHandler(db))

	// Notification routes (available to all authenticated users)
	r.GET("/notifications", handlers.AuthMiddleware(db), handlers.GetNotificationsHandler(db))
	r.GET("/notifications/unread-count", handlers.AuthMiddleware(db), handlers.GetUnreadCountHandler(db))
	r.PATCH("/notifications/:id/read", handlers.AuthMiddleware(db), handlers.MarkNotificationReadHandler(db))
	r.PATCH("/notifications/mark-all-read", handlers.AuthMiddleware(db), handlers.MarkAllNotificationsReadHandler(db))
	r.DELETE("/notifications/:id", handlers.AuthMiddleware(db), handlers.DeleteNotificationHandler(db))

	// Audit Log routes (Manager only)
	r.GET("/manager/audit-logs", handlers.AuthMiddleware(db), handlers.GetAuditLogsHandler(db))
	r.GET("/manager/audit-logs/stats", handlers.AuthMiddleware(db), handlers.GetAuditLogStatsHandler(db))
	r.GET("/manager/audit-logs/recent", handlers.AuthMiddleware(db), handlers.GetRecentAuditLogsHandler(db))
	r.GET("/manager/audit-logs/critical", handlers.AuthMiddleware(db), handlers.GetCriticalAuditLogsHandler(db))
	r.GET("/manager/audit-logs/failed", handlers.AuthMiddleware(db), handlers.GetFailedAuditLogsHandler(db))
	r.GET("/manager/audit-logs/:id", handlers.AuthMiddleware(db), handlers.GetAuditLogByIDHandler(db))
	r.GET("/manager/audit-logs/entity/:entity_type/:entity_id", handlers.AuthMiddleware(db), handlers.GetAuditLogsByEntityHandler(db))
	r.GET("/manager/audit-logs/actor/:actor_type/:actor_id", handlers.AuthMiddleware(db), handlers.GetAuditLogsByActorHandler(db))

	// Test endpoint for audit logging
	//r.GET("/test-audit-log", handlers.TestAuditLog(db))

	// n8n Webhook routes (API key authentication instead of JWT)
	n8n := r.Group("/n8n")
	n8n.Use(handlers.N8nAPIKeyMiddleware())
	{
		n8n.GET("/health", handlers.N8nHealthCheckHandler())

		// ⭐ PRIMARY: Process email with Gemini AI - receives prompt + id from n8n
		// This endpoint calls Gemini API internally and creates ticket
		n8n.POST("/process-email", handlers.ProcessEmailHandler(db))

		// Legacy: Direct ticket creation (IDs or names)
		n8n.POST("/ticket", handlers.N8nCreateTicketHandler(db))

		// Legacy: Smart ticket creation with pre-extracted AI hints
		n8n.POST("/smart-ticket", handlers.SmartTicketHandler(db))

		// Lookup endpoints (for manual testing or AI tools)
		n8n.GET("/lookup/accounts", handlers.N8nLookupAccountsHandler(db))
		n8n.GET("/lookup/contacts", handlers.N8nLookupContactsHandler(db))
		n8n.GET("/lookup/products", handlers.N8nLookupProductsHandler(db))

		// AI Agent Tool endpoints - designed for function calling
		// These return structured data for AI to make decisions
		tools := n8n.Group("/tools")
		{
			// Contact lookup tools
			tools.GET("/contact/by-email", handlers.ToolSearchContactByEmail(db))
			tools.GET("/contact/by-phone", handlers.ToolSearchContactByPhone(db))
			tools.GET("/contact/by-name", handlers.ToolSearchContactByName(db))

			// Account lookup tools
			tools.GET("/account/by-name", handlers.ToolSearchAccountByName(db))
			tools.GET("/account/by-domain", handlers.ToolSearchAccountByDomain(db))
			tools.GET("/account/contacts", handlers.ToolGetAccountContacts(db))

			// Product lookup tools
			tools.GET("/product/search", handlers.ToolSearchProduct(db))
			tools.GET("/products", handlers.ToolListAllProducts(db))
		}
	}

	return r
}
