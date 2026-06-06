package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type varianceRow struct {
	EntryID          uint    `json:"entryId"`
	PlanID           uint    `json:"planId"`
	PlanNumber       string  `json:"planNumber"`
	ProductID        uint    `json:"productId"`
	ProductName      string  `json:"productName"`
	ProductionDate   string  `json:"productionDate"`
	ProducedQuantity int     `json:"producedQuantity"`
	PlannedQuantity  int     `json:"plannedQuantity"`
	Difference       int     `json:"difference"`
	AchievementPct   float64 `json:"achievementPct"`
}

func VarianceReport(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rows []varianceRow

		sql := `SELECT
            pe.entry_id AS entry_id,
            pe.plan_id AS plan_id,
            pp.plan_number AS plan_number,
            p.product_id AS product_id,
            p.product_name AS product_name,
            TO_CHAR(pe.production_date, 'YYYY-MM-DD') AS production_date,
            pe.produced_quantity AS produced_quantity,
            pp.planned_quantity AS planned_quantity,
            (pe.produced_quantity - pp.planned_quantity) AS difference,
            CASE
                WHEN pp.planned_quantity = 0 THEN 0
                ELSE ROUND((pe.produced_quantity::numeric / pp.planned_quantity::numeric) * 100, 2)
            END AS achievement_pct
        FROM production_entries pe
        JOIN production_plans pp ON pe.plan_id = pp.plan_id
        JOIN products p ON pp.product_id = p.product_id
        ORDER BY pe.production_date DESC, pe.entry_id DESC;`

		if err := db.Raw(sql).Scan(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to generate variance report"})
			return
		}

		c.JSON(http.StatusOK, rows)
	}
}
