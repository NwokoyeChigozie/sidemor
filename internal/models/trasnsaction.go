package models

import (
	"fmt"
	"time"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	CustomerID      int64     `gorm:"column:customer_id; type:int" json:"customer_id"`
	PaymentModuleID int64     `gorm:"column:payment_module_id; type:int" json:"payment_module_id"`
	Reference       string    `gorm:"column:reference; type:varchar(255)" json:"reference"`
	Description     string    `gorm:"column:description; type:varchar(255)" json:"description"`
	Amount          float64   `gorm:"column:amount; type:decimal(20,2)" json:"amount"`
	PaymentMethod   string    `gorm:"column:payment_method; type:varchar(255); comment: (card, bank transfer, mobile money etc)" json:"payment_method"`
	Status          string    `gorm:"column:status; type:varchar(255)" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (t *Transaction) CreateTransaction(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("Transaction creation failed: %v", err.Error())
	}
	return nil
}
