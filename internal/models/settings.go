package models

import "time"

type Setting struct {
	ID             uint                   `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID      int64                  `gorm:"column:account_id; type:int; not null" json:"account_id"`
	BusinessTypeID int64                  `gorm:"column:business_type_id; type:int" json:"business_type_id"`
	UsageType      string                 `gorm:"column:usage_type; type:varchar(255)" json:"usage_type"`
	Countries      []SettingsCountries    `gorm:"column:countries;serializer:json" json:"countries"`
	Verifications  []SettingsVerification `gorm:"column:verifications;serializer:json" json:"verifications"`
	CurrencyCodes  []string               `gorm:"column:currency_codes;serializer:json" json:"currency_codes"`
	CreatedAt      time.Time              `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time              `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type SettingsCountries struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	CurrencyCode string `json:"currency_code"`
}

type SettingsVerification struct {
	DocumentUrl string             `json:"document_url"`
	Status      VerificationStatus `json:"status"`
	CountryID   uint               `json:"country_id"`
}
