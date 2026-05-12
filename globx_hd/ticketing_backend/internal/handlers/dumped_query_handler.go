package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DumpedQueryHandler struct {
	db *gorm.DB
}

func NewDumpedQueryHandler(db *gorm.DB) *DumpedQueryHandler {
	return &DumpedQueryHandler{db: db}
}

// GetDumpedQueries lists all dumped queries with pagination and filtering
func (h *DumpedQueryHandler) GetDumpedQueries(c *gin.Context) {
	var queries []models.DumpedQuery
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	offset := (page - 1) * limit

	query := h.db.Model(&models.DumpedQuery{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	result := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&queries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dumped queries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  queries,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetDumpedQuery fetches a single dumped query details
func (h *DumpedQueryHandler) GetDumpedQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.DumpedQuery

	if err := h.db.First(&query, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dumped query not found"})
		return
	}

	c.JSON(http.StatusOK, query)
}

// DeleteDumpedQuery deletes a dumped query
func (h *DumpedQueryHandler) DeleteDumpedQuery(c *gin.Context) {
	id := c.Param("id")

	// Get query first for audit logging
	var query models.DumpedQuery
	if err := h.db.First(&query, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dumped query not found"})
		return
	}

	if err := h.db.Delete(&models.DumpedQuery{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete dumped query"})
		return
	}

	// Audit log
	auditService := services.NewAuditService(h.db)
	auditService.LogCRUD(
		c,
		models.AuditDumpedQueryDeleted,
		models.EntityTypeDumpedQuery,
		&query.ID,
		query.Subject,
		fmt.Sprintf("Dumped query deleted: %s", query.Subject),
		query,
		nil,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Dumped query deleted successfully"})
}

// UpdateDumpedQueryStatus updates the status (e.g., RESOLVED, IGNORED)
func (h *DumpedQueryHandler) UpdateDumpedQueryStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get query first for audit logging
	var query models.DumpedQuery
	if err := h.db.First(&query, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dumped query not found"})
		return
	}
	oldStatus := query.Status

	if err := h.db.Model(&models.DumpedQuery{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	// Audit log
	auditService := services.NewAuditService(h.db)
	auditService.LogCRUD(
		c,
		models.AuditDumpedQueryResolved,
		models.EntityTypeDumpedQuery,
		&query.ID,
		query.Subject,
		fmt.Sprintf("Dumped query status changed: %s -> %s", oldStatus, req.Status),
		map[string]string{"status": oldStatus},
		map[string]string{"status": req.Status},
	)

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
