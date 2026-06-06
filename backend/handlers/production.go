package handlers

import (
	"net/http"
	"strconv"
	"time"

	"projectplanningtracking/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productionEntryInput struct {
	PlanID           uint   `json:"planId" binding:"required"`
	ProductionDate   string `json:"productionDate" binding:"required"`
	ProducedQuantity int    `json:"producedQuantity" binding:"required"`
	Shift            string `json:"shift"`
	Remarks          string `json:"remarks"`
}

func ListProductionEntries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := getUserContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var entries []models.ProductionEntry
		query := db.Preload("Plan.Product")
		if role == "operator" {
			query = query.Joins("Plan").Where("production_plans.assigned_to = ?", userID)
		}

		if err := query.Find(&entries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch production entries"})
			return
		}
		c.JSON(http.StatusOK, entries)
	}
}

func CreateProductionEntry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := getUserContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var input productionEntryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var plan models.ProductionPlan
		if err := db.Preload("Product").First(&plan, input.PlanID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan selected"})
			return
		}

		if role == "operator" {
			if plan.AssignedTo == nil || *plan.AssignedTo != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized for selected plan"})
				return
			}
		}

		productionDate, err := time.Parse("2006-01-02", input.ProductionDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "productionDate must be in YYYY-MM-DD format"})
			return
		}

		entry := models.ProductionEntry{
			PlanID:           input.PlanID,
			ProductionDate:   productionDate,
			ProducedQuantity: input.ProducedQuantity,
			Shift:            input.Shift,
			Remarks:          input.Remarks,
		}

		if err := db.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create production entry"})
			return
		}

		if err := db.Preload("Plan.Product").First(&entry, entry.EntryID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch created entry"})
			return
		}

		c.JSON(http.StatusCreated, entry)
	}
}

func UpdateProductionEntry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := getUserContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var entry models.ProductionEntry
		if err := db.Preload("Plan").First(&entry, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "production entry not found"})
			return
		}

		if role == "operator" {
			if entry.Plan.AssignedTo == nil || *entry.Plan.AssignedTo != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to modify this entry"})
				return
			}
		}

		var input productionEntryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.PlanID != entry.PlanID {
			var plan models.ProductionPlan
			if err := db.First(&plan, input.PlanID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan selected"})
				return
			}
			if role == "operator" {
				if plan.AssignedTo == nil || *plan.AssignedTo != userID {
					c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized for selected plan"})
					return
				}
			}
		}

		productionDate, err := time.Parse("2006-01-02", input.ProductionDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "productionDate must be in YYYY-MM-DD format"})
			return
		}

		updates := map[string]interface{}{
			"plan_id":           input.PlanID,
			"production_date":   productionDate,
			"produced_quantity": input.ProducedQuantity,
			"shift":             input.Shift,
			"remarks":           input.Remarks,
		}

		if err := db.Model(&entry).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update production entry"})
			return
		}

		if err := db.Preload("Plan.Product").First(&entry, entry.EntryID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch updated entry"})
			return
		}

		c.JSON(http.StatusOK, entry)
	}
}
