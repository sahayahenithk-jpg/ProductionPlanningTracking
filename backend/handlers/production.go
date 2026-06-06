package handlers

import (
	"net/http"
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
		var entries []models.ProductionEntry
		if err := db.Preload("Plan.Product").Find(&entries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch production entries"})
			return
		}
		c.JSON(http.StatusOK, entries)
	}
}

func CreateProductionEntry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
