package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionStatus string

var (
	TransactionSuccessful TransactionStatus = "successful"
	TransactionPending    TransactionStatus = "pending"
	TransactionFailed     TransactionStatus = "failed"
)

type Transaction struct {
	ID              uint              `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	MerchantID      int64             `gorm:"column:merchant_id; type:int" json:"merchant_id"`
	CustomerID      int64             `gorm:"column:customer_id; type:int" json:"customer_id"`
	CustomerName    string            `gorm:"-" json:"customer_name"`
	PaymentModuleID int64             `gorm:"column:payment_module_id; type:int" json:"payment_module_id"`
	Reference       string            `gorm:"column:reference; type:varchar(255)" json:"reference"`
	MerchantName    string            `gorm:"-" json:"merchant_name"`
	MerchantEmail   string            `gorm:"-" json:"merchant_email"`
	Country         string            `gorm:"-" json:"country"`
	Currency        string            `gorm:"-" json:"currency"`
	Description     string            `gorm:"column:description; type:varchar(255)" json:"description"`
	Amount          float64           `gorm:"column:amount; type:decimal(20,2)" json:"amount"`
	TaxFee          float64           `gorm:"column:tax_fee; type:decimal(20,2)" json:"tax_fee"`
	ProcessingFee   float64           `gorm:"column:processing_fee; type:decimal(20,2)" json:"processing_fee"`
	CountryID       int64             `gorm:"column:country_id; type:int" json:"country_id"`
	PaymentMethod   PaymentMethod     `gorm:"column:payment_method; type:varchar(255); comment: (card, bank transfer, mobile money etc)" json:"payment_method"`
	Status          TransactionStatus `gorm:"column:status; type:varchar(255)" json:"status"`
	IsPaidOut       bool              `gorm:"column:is_paid_out; default: false" json:"is_paid_out"`
	PayoutID        int64             `gorm:"column:payout_id; type:int" json:"payout_id"`
	TransactionDate time.Time         `gorm:"column:transaction_date" json:"transaction_date"`
	CreatedAt       time.Time         `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type TransactionSummary struct {
	Currency  string  `gorm:"-" json:"currency"`
	CountryID int64   `gorm:"column:country_id; type:int" json:"-"`
	Amount    float64 `gorm:"column:amount; type:decimal(20,2)" json:"amount"`
}

type RecordTransactionRequest struct {
	AccountID            int64   `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	Reference            string  `json:"reference" validate:"required"`
	Description          string  `json:"description"`
	Country              int     `json:"country"  validate:"required"`
	Amount               float64 `json:"amount" validate:"required"`
	TaxFee               float64 `json:"tax_fee"`
	ProcessingFee        float64 `json:"processing_fee"`
	TransactionCreatedAt int     `json:"transaction_created_at"`
}

type GetTransactionsRequest struct {
	Search         string `json:"search"`
	CurrencyFilter string `json:"currency"`
	Status         string `json:"status"`
	FromTime       int    `json:"from_time"`
	ToTime         int    `json:"to_time"`
}

func (t *Transaction) GetTransactionByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "id = ?", t.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (t *Transaction) GetTransactionsSummary(db *gorm.DB, paidOut *bool) ([]TransactionSummary, error) {
	summary := []TransactionSummary{}
	extraQuery := ""
	whereQuery := ""
	if t.MerchantID != 0 {
		extraQuery += fmt.Sprintf(" and merchant_id = %v", t.MerchantID)
		whereQuery = addQuery(whereQuery, fmt.Sprintf(" and merchant_id = %v", t.MerchantID), "and")
	}

	if paidOut != nil {
		extraQuery += fmt.Sprintf(" and is_paid_out = %v", *paidOut)
		whereQuery = addQuery(whereQuery, fmt.Sprintf(" and is_paid_out = %v", *paidOut), "and")
	}

	// subQuery := db.Model(&t).Select("SUM(amount)").Where("is_paid_out = ? and currency=transactions.currency", false)
	selectQuery := fmt.Sprintf("country_id, (SUM(amount) from transactions where currency=transactions.currency " + extraQuery + ") as amount")

	_, err := postgresql.RawSelectAllFromByGroup(db, "country_id", "desc", nil, &t, &summary, "country_id", selectQuery, whereQuery)
	if err != nil {
		return summary, err
	}

	return summary, nil
}

func (t *Transaction) GetTransactions(db *gorm.DB, paginator postgresql.Pagination, search string, from int, to int, paidOut *bool) ([]Transaction, postgresql.PaginationResponse, error) {
	details := []Transaction{}
	query := ""

	if paidOut != nil {
		query = addQuery(query, fmt.Sprintf("is_paid_out = %v", *paidOut), "and")
	}

	if t.MerchantID != 0 {
		query = addQuery(query, fmt.Sprintf("merchant_id = %v", t.MerchantID), "and")
	}

	if search != "" {
		query = addQuery(query, fmt.Sprintf("reference = '%v'", search), "and")
	}

	if t.CountryID != 0 {
		query = addQuery(query, fmt.Sprintf("country_id = %v", t.CountryID), "and")
	}

	if t.Status != "" {
		query = addQuery(query, fmt.Sprintf("status = '%v'", t.Status), "and")
	}

	if from != 0 {
		fromTime := time.Unix(int64(from), 0)
		query = addQuery(query, fmt.Sprintf("transaction_date >= '%v'", fromTime), "and")
	}

	if to != 0 {
		toTime := time.Unix(int64(to), 0)
		query = addQuery(query, fmt.Sprintf("transaction_date >= '%v'", toTime), "and")
	}

	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, query)
	if err != nil {
		return details, pagination, err
	}

	return details, pagination, nil
}

func (t *Transaction) CreateTransaction(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("Transaction creation failed: %v", err.Error())
	}
	return nil
}

func (t *Transaction) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &t)
	return err
}
