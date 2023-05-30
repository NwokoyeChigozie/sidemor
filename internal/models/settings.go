package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type VerificationStatus string

var (
	NotVerified         VerificationStatus = "not_verified"
	VerificationPending VerificationStatus = "pending"
	Verified            VerificationStatus = "verified"
)

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

type SaveSettingsRequest struct {
	Countries           []int                         `json:"document_url"  validate:"required"`
	WalletCurrencyCodes []string                      `json:"wallet_currency_codes"  validate:"required"`
	BusinessTypeID      int64                         `json:"business_type_id"  validate:"required"`
	UsageType           string                        `json:"usage_type"  validate:"required,oneof=second minute hour day week month year"`
	IntervalBase        string                        `json:"interval_base" validate:"required,oneof=online offline"`
	Documents           []SettingsVerificationRequest `json:"documents" validate:"required"`
}

type SettingsVerificationRequest struct {
	DocumentUrl string `json:"document_url"`
	CountryID   uint   `json:"country_id"`
}

func (s *Setting) CreateSetting(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &s)
	if err != nil {
		return fmt.Errorf("setting creation failed: %v", err.Error())
	}
	return nil
}

func (s *Setting) GetSettingByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &s, "account_id = ?", s.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *Setting) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &s)
	return err
}
