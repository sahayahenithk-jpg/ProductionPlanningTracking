package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"

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

type summaryRow struct {
    TotalProducts   int     `json:"totalProducts"`
    TotalPlans      int     `json:"totalPlans"`
    ProductionTotal int     `json:"productionTotal"`
    AchievementPct  float64 `json:"achievementPct"`
}

func parseDateParam(value string) (time.Time, error) {
    return time.Parse("2006-01-02", value)
}

func parseMonthParam(value string) (time.Time, time.Time, error) {
    start, err := time.Parse("2006-01", value)
    if err != nil {
        return time.Time{}, time.Time{}, err
    }
    end := start.AddDate(0, 1, 0)
    return start, end, nil
}

func buildVarianceFilter(c *gin.Context) (string, []interface{}, error) {
    filters := []string{"1=1"}
    params := []interface{}{}

    if productIDStr := c.Query("productId"); productIDStr != "" {
        productID, err := strconv.Atoi(productIDStr)
        if err != nil {
            return "", nil, fmt.Errorf("invalid productId")
        }
        filters = append(filters, "pp.product_id = ?")
        params = append(params, productID)
    }

    if month := c.Query("month"); month != "" {
        start, end, err := parseMonthParam(month)
        if err != nil {
            return "", nil, err
        }
        filters = append(filters, "pe.production_date >= ?", "pe.production_date < ?")
        params = append(params, start.Format("2006-01-02"), end.Format("2006-01-02"))
    } else {
        if startDate := c.Query("startDate"); startDate != "" {
            start, err := parseDateParam(startDate)
            if err != nil {
                return "", nil, err
            }
            filters = append(filters, "pe.production_date >= ?")
            params = append(params, start.Format("2006-01-02"))
        }
        if endDate := c.Query("endDate"); endDate != "" {
            end, err := parseDateParam(endDate)
            if err != nil {
                return "", nil, err
            }
            filters = append(filters, "pe.production_date <= ?")
            params = append(params, end.Format("2006-01-02"))
        }
    }

    return strings.Join(filters, " AND "), params, nil
}

func VarianceReport(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var rows []varianceRow
        filterClause, params, err := buildVarianceFilter(c)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        sql := fmt.Sprintf(`SELECT
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
        WHERE %s
        ORDER BY pe.production_date DESC, pe.entry_id DESC;`, filterClause)

        if err := db.Raw(sql, params...).Scan(&rows).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to generate variance report"})
            return
        }

        c.JSON(http.StatusOK, rows)
    }
}

func buildSummaryClauses(c *gin.Context) (string, []interface{}, string, []interface{}, interface{}, error) {
    planFilters := []string{"1=1"}
    planParams := []interface{}{}
    entryFilters := []string{"1=1"}
    entryParams := []interface{}{}
    var productID interface{}

    if productIDStr := c.Query("productId"); productIDStr != "" {
        id, err := strconv.Atoi(productIDStr)
        if err != nil {
            return "", nil, "", nil, nil, fmt.Errorf("invalid productId")
        }
        productID = id
        planFilters = append(planFilters, "product_id = ?")
        planParams = append(planParams, id)
        entryFilters = append(entryFilters, "pp.product_id = ?")
        entryParams = append(entryParams, id)
    }

    if month := c.Query("month"); month != "" {
        start, end, err := parseMonthParam(month)
        if err != nil {
            return "", nil, "", nil, nil, err
        }
        planFilters = append(planFilters, "plan_date >= ?", "plan_date < ?")
        planParams = append(planParams, start.Format("2006-01-02"), end.Format("2006-01-02"))
        entryFilters = append(entryFilters, "pe.production_date >= ?", "pe.production_date < ?")
        entryParams = append(entryParams, start.Format("2006-01-02"), end.Format("2006-01-02"))
    } else {
        if startDate := c.Query("startDate"); startDate != "" {
            start, err := parseDateParam(startDate)
            if err != nil {
                return "", nil, "", nil, nil, err
            }
            planFilters = append(planFilters, "plan_date >= ?")
            planParams = append(planParams, start.Format("2006-01-02"))
            entryFilters = append(entryFilters, "pe.production_date >= ?")
            entryParams = append(entryParams, start.Format("2006-01-02"))
        }
        if endDate := c.Query("endDate"); endDate != "" {
            end, err := parseDateParam(endDate)
            if err != nil {
                return "", nil, "", nil, nil, err
            }
            planFilters = append(planFilters, "plan_date <= ?")
            planParams = append(planParams, end.Format("2006-01-02"))
            entryFilters = append(entryFilters, "pe.production_date <= ?")
            entryParams = append(entryParams, end.Format("2006-01-02"))
        }
    }

    return strings.Join(planFilters, " AND "), planParams, strings.Join(entryFilters, " AND "), entryParams, productID, nil
}

func SummaryReport(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var summary summaryRow
        planClause, planParams, entryClause, entryParams, productID, err := buildSummaryClauses(c)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if productID != nil {
            if err := db.Raw("SELECT COUNT(*) FROM products WHERE product_id = ?", productID).Scan(&summary.TotalProducts).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to compute total products"})
                return
            }
        } else {
            if err := db.Raw("SELECT COUNT(*) FROM products").Scan(&summary.TotalProducts).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to compute total products"})
                return
            }
        }

        if err := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM production_plans WHERE %s", planClause), planParams...).Scan(&summary.TotalPlans).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to compute total plans"})
            return
        }

        if err := db.Raw(fmt.Sprintf(`SELECT COALESCE(SUM(pe.produced_quantity), 0)
            FROM production_entries pe
            JOIN production_plans pp ON pe.plan_id = pp.plan_id
            WHERE %s`, entryClause), entryParams...).Scan(&summary.ProductionTotal).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to compute production total"})
            return
        }

        var producedTotal int64
        var plannedTotal int64
        if err := db.Raw(fmt.Sprintf(`SELECT COALESCE(SUM(pe.produced_quantity), 0), COALESCE(SUM(pp.planned_quantity), 0)
            FROM production_entries pe
            JOIN production_plans pp ON pe.plan_id = pp.plan_id
            WHERE %s`, entryClause), entryParams...).Row().Scan(&producedTotal, &plannedTotal); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to compute achievement"})
            return
        }

        if plannedTotal > 0 {
            summary.AchievementPct = float64(producedTotal) / float64(plannedTotal) * 100
        }

        c.JSON(http.StatusOK, summary)
    }
}
