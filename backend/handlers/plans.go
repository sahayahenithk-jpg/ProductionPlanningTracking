package handlers

import (
	"net/http"
	"strconv"
	"time"

	"projectplanningtracking/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type planInput struct {
	PlanNumber      string `json:"planNumber" binding:"required"`
	ProductID       uint   `json:"productId" binding:"required"`
	PlanDate        string `json:"planDate" binding:"required"`
	PlannedQuantity int    `json:"plannedQuantity" binding:"required"`
	AssignedTo      *uint  `json:"assignedTo"`
	Remarks         string `json:"remarks"`
	Status          string `json:"status"`
}

func parsePlanDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func ListPlans(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := getUserContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var plans []models.ProductionPlan
		query := db.Preload("Product").Preload("AssignedUser")
		if role == "operator" {
			query = query.Where("production_plans.assigned_to = ?", userID)
		}

		if err := query.Find(&plans).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch plans"})
			return
		}
		c.JSON(http.StatusOK, plans)
	}
}

func CreatePlan(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input planInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.First(&models.Product{}, input.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product selected"})
			return
		}

		if input.AssignedTo != nil {
			var user models.User
			if err := db.First(&user, *input.AssignedTo).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assigned user"})
				return
			}
		}

		planDate, err := parsePlanDate(input.PlanDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "planDate must be in YYYY-MM-DD format"})
			return
		}

		plan := models.ProductionPlan{
			PlanNumber:      input.PlanNumber,
			ProductID:       input.ProductID,
			PlanDate:        planDate,
			PlannedQuantity: input.PlannedQuantity,
			AssignedTo:      input.AssignedTo,
			Remarks:         input.Remarks,
			Status:          input.Status,
		}

		if err := db.Create(&plan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create plan"})
			return
		}

		if err := db.Preload("Product").Preload("AssignedUser").First(&plan, plan.PlanID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch created plan"})
			return
		}

		c.JSON(http.StatusCreated, plan)
	}
}

func UpdatePlan(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var plan models.ProductionPlan
		if err := db.First(&plan, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}

		var input planInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.First(&models.Product{}, input.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product selected"})
			return
		}

		if input.AssignedTo != nil {
			var user models.User
			if err := db.First(&user, *input.AssignedTo).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assigned user"})
				return
			}
		}

		planDate, err := parsePlanDate(input.PlanDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "planDate must be in YYYY-MM-DD format"})
			return
		}

		updates := map[string]interface{}{
			"plan_number":      input.PlanNumber,
			"product_id":       input.ProductID,
			"plan_date":        planDate,
			"planned_quantity": input.PlannedQuantity,
			"assigned_to":      input.AssignedTo,
			"remarks":          input.Remarks,
			"status":           input.Status,
		}

		if err := db.Model(&plan).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update plan"})
			return
		}

		if err := db.Preload("Product").Preload("AssignedUser").First(&plan, plan.PlanID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch updated plan"})
			return
		}

		c.JSON(http.StatusOK, plan)
	}
}
