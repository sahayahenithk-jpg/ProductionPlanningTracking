package models

import "time"

type ProductionPlan struct {
	PlanID          uint      `gorm:"primaryKey" json:"planId"`
	PlanNumber      string    `gorm:"size:100;unique;not null" json:"planNumber"`
	ProductID       uint      `gorm:"not null" json:"productId"`
	Product         Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product"`
	PlanDate        time.Time `json:"planDate"`
	PlannedQuantity int       `json:"plannedQuantity"`
	Remarks         string    `gorm:"type:text" json:"remarks"`
	Status          string    `gorm:"size:50" json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}
