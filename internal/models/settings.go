package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type VerificationStatus string
type PaymentMethod string

var (
	NotVerified         VerificationStatus = "not_verified"
	VerificationPending VerificationStatus = "pending"
	Verified            VerificationStatus = "verified"
)

var (
	CardMethod              PaymentMethod = "card"
	AccountMethod           PaymentMethod = "account"
	BankTransferMethod      PaymentMethod = "banktransfer"
	MpesaMethod             PaymentMethod = "mpesa"
	MobileMoneyGhanaMethod  PaymentMethod = "mobilemoneyghana"
	MobileMoneyFrancoMethod PaymentMethod = "mobilemoneyfranco"
	MobileMoneyUgandaMethod PaymentMethod = "mobilemoneyuganda"
	MobileMoneyRwandaMethod PaymentMethod = "mobilemoneyrwanda"
	MobileMoneyZambiaMethod PaymentMethod = "mobilemoneyzambia"
	BarterMethod            PaymentMethod = "barter"
	NqrMethod               PaymentMethod = "nqr"
	UssdMethod              PaymentMethod = "ussd"
	CreditMethod            PaymentMethod = "credit"
)

type Setting struct {
	ID             uint                   `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID      int64                  `gorm:"column:account_id; type:int; not null" json:"account_id"`
	BusinessTypeID int64                  `gorm:"column:business_type_id; type:int" json:"business_type_id"`
	UsageType      string                 `gorm:"column:usage_type; type:varchar(255)" json:"usage_type"`
	Countries      []SettingsCountries    `gorm:"column:countries;serializer:json" json:"countries"`
	Verifications  []SettingsVerification `gorm:"column:verifications;serializer:json" json:"verifications"`
	CurrencyCodes  []string               `gorm:"column:currency_codes;serializer:json" json:"currency_codes"`
	PaymentMethods []PaymentMethod        `gorm:"column:payment_methods;serializer:json" json:"payment_methods"`
	IsVerified     bool                   `gorm:"column:is_verified; default:false" json:"is_verified"`
	AccountType    string                 `gorm:"-" json:"account_type"`
	Email          string                 `gorm:"-" json:"email"`
	FullName       string                 `gorm:"-" json:"full_name"`
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
	Countries           []int                         `json:"countries"  validate:"required"`
	WalletCurrencyCodes []string                      `json:"wallet_currency_codes"  validate:"required"`
	BusinessTypeID      int64                         `json:"business_type_id"  validate:"required"`
	UsageType           string                        `json:"usage_type"  validate:"required,oneof=online offline"`
	Documents           []SettingsVerificationRequest `json:"documents" validate:"required"`
}
type EnableOrDisablePaymentMethodsRequest struct {
	Methods []PaymentMethod `json:"methods"  validate:"required"`
}

type AddRemoveOrGetWalletsRequest struct {
	CurrencyCodes []string `json:"currency_codes"`
}

type SettingsVerificationRequest struct {
	DocumentUrl string `json:"document_url"`
	CountryID   uint   `json:"country_id"`
}

type GetSettingsRequest struct {
	Search   string `json:"search"`
	Status   string `json:"status"`
	FromTime int    `json:"from_time"`
	ToTime   int    `json:"to_time"`
}

type UpdateDocumentStatusRequest struct {
	CountryId int    `json:"country_id"`
	Status    string `json:"status" validate:"oneof=not_verified pending verified"`
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

func (s *Setting) GetSettingByID(db *gorm.DB, settingsID int) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &s, "id = ?", settingsID)
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

func (p PaymentMethod) In(methods []PaymentMethod) bool {
	for _, v := range methods {
		if p == v {
			return true
		}
	}
	return false
}

func (s *Setting) GetSettings(db *gorm.DB, paginator postgresql.Pagination, userIds []int, from int, to int, isVerified *bool) ([]Setting, postgresql.PaginationResponse, error) {
	details := []Setting{}
	query := ""

	if isVerified != nil {
		query = addQuery(query, fmt.Sprintf("is_verified = %v", *isVerified), "and")
	}

	if len(userIds) > 0 {
		idsString := []string{}
		for _, v := range userIds {
			idsString = append(idsString, fmt.Sprintf("%v", v))
		}
		query = addQuery(query, fmt.Sprintf(" account_id IN (%s)", strings.Join(idsString, ",")), "and")
	}

	if from != 0 {
		fromTime := time.Unix(int64(from), 0)
		query = addQuery(query, fmt.Sprintf("created_at >= '%v'", fromTime), "and")
	}

	if to != 0 {
		toTime := time.Unix(int64(to), 0)
		query = addQuery(query, fmt.Sprintf("created_at >= '%v'", toTime), "and")
	}

	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, query)
	if err != nil {
		return details, pagination, err
	}

	return details, pagination, nil
}

func (s *Setting) UpdateVerificationSettings(db *gorm.DB, request UpdateDocumentStatusRequest) error {
	var err error

	found := false
	for i, v := range s.Verifications {
		if v.CountryID == uint(request.CountryId) {
			s.Verifications[i].Status = VerificationStatus(request.Status)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("country ID not found")
	}

	err = db.Save(&s).Error
	if err != nil {
		return err
	}
	return err
}
