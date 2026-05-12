package handlers

import (
	"fmt"
	"net/http"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Products
type CreateProductInput struct {
	ProductName        string `json:"product_name" binding:"required"`
	ProductDescription string `json:"product_description"`
}

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateProductInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := models.MasterProduct{
			ProductName:        in.ProductName,
			ProductDescription: in.ProductDescription,
		}
		if err := db.Create(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditProductCreated, models.EntityTypeProduct, &p.ID, p.ProductName, fmt.Sprintf("Product created: %s", p.ProductName), nil, p)

		c.JSON(http.StatusCreated, p)
	}
}

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.MasterProduct
		if err := db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func GetProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var p models.MasterProduct
		if err := db.First(&p, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	}
}

type UpdateProductInput struct {
	ProductName        *string `json:"product_name"`
	ProductDescription *string `json:"product_description"`
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var p models.MasterProduct
		if err := db.First(&p, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		oldProduct := p
		var in UpdateProductInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.ProductName != nil {
			p.ProductName = *in.ProductName
		}
		if in.ProductDescription != nil {
			p.ProductDescription = *in.ProductDescription
		}
		if err := db.Save(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditProductUpdated, models.EntityTypeProduct, &p.ID, p.ProductName, fmt.Sprintf("Product updated: %s", p.ProductName), oldProduct, p)

		c.JSON(http.StatusOK, p)
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var p models.MasterProduct
		if err := db.First(&p, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err := db.Delete(&models.MasterProduct{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditProductDeleted, models.EntityTypeProduct, &p.ID, p.ProductName, fmt.Sprintf("Product deleted: %s", p.ProductName), p, nil)

		c.Status(http.StatusNoContent)
	}
}

// --- Roles & Designations
type CreateRoleInput struct {
	RoleName string `json:"role_name" binding:"required"`
}

// --- Issues
type CreateIssueInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	IssueName string `json:"issue_name" binding:"required"`
}

func CreateIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateIssueInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// ensure product exists
		var p models.MasterProduct
		if err := db.First(&p, in.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}
		item := models.MasterProductIssue{ProductID: in.ProductID, IssueName: in.IssueName}
		if err := db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditProductIssueCreated, models.EntityTypeProductIssue, &item.ID, item.IssueName, fmt.Sprintf("Product issue created: %s", item.IssueName), nil, item)

		c.JSON(http.StatusCreated, item)
	}
}

func GetIssues(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.MasterProductIssue
		// optional filter by product_id
		if pid := c.Query("product_id"); pid != "" {
			if err := db.Where("product_id = ?", pid).Find(&list).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, list)
			return
		}
		if err := db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

type UpdateIssueInput struct {
	ProductID *uint   `json:"product_id"`
	IssueName *string `json:"issue_name"`
}

func UpdateIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var item models.MasterProductIssue
		if err := db.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
			return
		}
		oldIssue := item
		var in UpdateIssueInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.ProductID != nil {
			item.ProductID = *in.ProductID
		}
		if in.IssueName != nil {
			item.IssueName = *in.IssueName
		}
		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditProductIssueUpdated, models.EntityTypeProductIssue, &item.ID, item.IssueName, fmt.Sprintf("Product issue updated: %s", item.IssueName), oldIssue, item)

		c.JSON(http.StatusOK, item)
	}
}

func DeleteIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var item models.MasterProductIssue
		if err := db.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
			return
		}
		if err := db.Delete(&models.MasterProductIssue{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditProductIssueDeleted, models.EntityTypeProductIssue, &item.ID, item.IssueName, fmt.Sprintf("Product issue deleted: %s", item.IssueName), item, nil)

		c.Status(http.StatusNoContent)
	}
}

func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateRoleInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		r := models.MasterRole{RoleName: in.RoleName}
		if err := db.Create(&r).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditRoleCreated, models.EntityTypeRole, &r.ID, r.RoleName, fmt.Sprintf("Role created: %s", r.RoleName), nil, r)

		c.JSON(http.StatusCreated, r)
	}
}

func GetRoles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.MasterRole
		if err := db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

type CreateDesignationInput struct {
	DesignationName string `json:"designation_name" binding:"required"`
}

func CreateUserDesignation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateDesignationInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		d := models.MasterUserDesignation{DesignationName: in.DesignationName}
		if err := db.Create(&d).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditDesignationCreated, models.EntityTypeDesignation, &d.ID, d.DesignationName, fmt.Sprintf("User designation created: %s", d.DesignationName), nil, d)

		c.JSON(http.StatusCreated, d)
	}
}

type UpdateRoleInput struct {
	RoleName *string `json:"role_name"`
}

func UpdateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var r models.MasterRole
		if err := db.First(&r, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		oldRole := r
		var in UpdateRoleInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.RoleName != nil {
			r.RoleName = *in.RoleName
		}
		if err := db.Save(&r).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditRoleUpdated, models.EntityTypeRole, &r.ID, r.RoleName, fmt.Sprintf("Role updated: %s", r.RoleName), oldRole, r)

		c.JSON(http.StatusOK, r)
	}
}

func DeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var r models.MasterRole
		if err := db.First(&r, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		if err := db.Delete(&models.MasterRole{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditRoleDeleted, models.EntityTypeRole, &r.ID, r.RoleName, fmt.Sprintf("Role deleted: %s", r.RoleName), r, nil)

		c.Status(http.StatusNoContent)
	}
}

func CreateContactDesignation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateDesignationInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		d := models.MasterContactDesignation{DesignationName: in.DesignationName}
		if err := db.Create(&d).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditDesignationCreated, models.EntityTypeDesignation, &d.ID, d.DesignationName, fmt.Sprintf("Contact designation created: %s", d.DesignationName), nil, d)

		c.JSON(http.StatusCreated, d)
	}
}

func GetUserDesignations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.MasterUserDesignation
		if err := db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func GetContactDesignations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.MasterContactDesignation
		if err := db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

type UpdateDesignationInput struct {
	DesignationName *string `json:"designation_name"`
}

func UpdateUserDesignation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var d models.MasterUserDesignation
		if err := db.First(&d, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "designation not found"})
			return
		}
		oldDesignation := d
		var in UpdateDesignationInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.DesignationName != nil {
			d.DesignationName = *in.DesignationName
		}
		if err := db.Save(&d).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditDesignationUpdated, models.EntityTypeDesignation, &d.ID, d.DesignationName, fmt.Sprintf("User designation updated: %s", d.DesignationName), oldDesignation, d)

		c.JSON(http.StatusOK, d)
	}
}

func DeleteUserDesignation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var d models.MasterUserDesignation
		if err := db.First(&d, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "designation not found"})
			return
		}
		if err := db.Delete(&models.MasterUserDesignation{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditDesignationDeleted, models.EntityTypeDesignation, &d.ID, d.DesignationName, fmt.Sprintf("User designation deleted: %s", d.DesignationName), d, nil)

		c.Status(http.StatusNoContent)
	}
}

func UpdateContactDesignation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var d models.MasterContactDesignation
		if err := db.First(&d, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "designation not found"})
			return
		}
		oldDesignation := d
		var in UpdateDesignationInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.DesignationName != nil {
			d.DesignationName = *in.DesignationName
		}
		if err := db.Save(&d).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditDesignationUpdated, models.EntityTypeDesignation, &d.ID, d.DesignationName, fmt.Sprintf("Contact designation updated: %s", d.DesignationName), oldDesignation, d)

		c.JSON(http.StatusOK, d)
	}
}

func DeleteContactDesignation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var d models.MasterContactDesignation
		if err := db.First(&d, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "designation not found"})
			return
		}
		if err := db.Delete(&models.MasterContactDesignation{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditDesignationDeleted, models.EntityTypeDesignation, &d.ID, d.DesignationName, fmt.Sprintf("Contact designation deleted: %s", d.DesignationName), d, nil)

		c.Status(http.StatusNoContent)
	}
}
