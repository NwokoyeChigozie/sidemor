package models

import "time"

type VerificationStatus string

var (
	NotVerified VerificationStatus = "not_verified"
	Pending     VerificationStatus = "pending"
	Verified    VerificationStatus = "verified"
)

type PaymentModule struct {
	ID               uint                        `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID        int64                       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	LogoUrl          string                      `gorm:"column:logo_url; type:varchar(255)" json:"logo_url"`
	Name             string                      `gorm:"column:name; type:varchar(255)" json:"name"`
	BackgroundColour string                      `gorm:"column:background_colour; type:varchar(255)" json:"background_colour"`
	ButtonColour     string                      `gorm:"column:button_colour; type:varchar(255)" json:"button_colour"`
	CountryID        int64                       `gorm:"column:country_id; type:int" json:"country_id"`
	CurrencyCode     string                      `gorm:"column:currency_code; type:varchar(255)" json:"currency_code"`
	IsShippingType   bool                        `gorm:"column:is_shipping_type; default: false" json:"is_shipping_type"`
	IsPublished      bool                        `gorm:"column:is_published; default: false" json:"is_published"`
	Vat              float64                     `gorm:"column:vat; type:decimal(20,2)" json:"vat"`
	ShippingTypes    []PaymentModuleShippingType `gorm:"column:shipping_types;serializer:json" json:"shipping_types"`
	CreatedAt        time.Time                   `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time                   `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type PaymentModuleShippingType struct {
	Name         string  `json:"name"`
	Time         string  `json:"time"`
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currency_code"`
}
