package models

import "time"

type ProductionEntry struct {
	EntryID          uint           `gorm:"primaryKey" json:"entryId"`
	PlanID           uint           `gorm:"not null" json:"planId"`
	Plan             ProductionPlan `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"plan"`
	ProductionDate   time.Time      `json:"productionDate"`
	ProducedQuantity int            `json:"producedQuantity"`
	Shift            string         `gorm:"size:50" json:"shift"`
	Remarks          string         `gorm:"type:text" json:"remarks"`
	CreatedAt        time.Time      `json:"createdAt"`
}
