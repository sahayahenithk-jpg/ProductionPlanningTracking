package models

import "time"

type Product struct {
	ProductID   uint      `gorm:"primaryKey" json:"productId"`
	ProductCode string    `gorm:"size:100;unique;not null" json:"productCode"`
	ProductName string    `gorm:"size:200;not null" json:"productName"`
	Unit        string    `gorm:"size:50" json:"unit"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:50" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}
