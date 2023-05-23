package models

import "time"

type PaymentOrder struct {
	ID               uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	CustomerID       int64     `gorm:"column:customer_id; type:int" json:"customer_id"`
	PaymentHistoryID int64     `gorm:"column:payment_history_id; type:int" json:"payment_history_id"`
	Item             string    `gorm:"column:item; type:varchar(255)" json:"item"`
	Quantity         int64     `gorm:"column:quantity; type:int; default:1" json:"quantity"`
	UnitPrice        float64   `gorm:"column:unit_price; type:decimal(20,2)" json:"unit_price"`
	CreatedAt        time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}
