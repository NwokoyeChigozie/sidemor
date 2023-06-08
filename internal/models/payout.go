package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Payout struct {
	ID            uint              `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	MerchantID    int64             `gorm:"column:merchant_id; type:int" json:"merchant_id"`
	Reference     string            `gorm:"column:reference; type:varchar(255)" json:"reference"`
	MerchantName  string            `gorm:"-" json:"merchant_name"`
	MerchantEmail string            `gorm:"-" json:"merchant_email"`
	Currency      string            `gorm:"-" json:"Currency"`
	Amount        float64           `gorm:"column:amount; type:decimal(20,2)" json:"amount"`
	CountryID     int64             `gorm:"column:country_id; type:int" json:"country_id"`
	Status        TransactionStatus `gorm:"column:status; type:varchar(255)" json:"status"`
	CreatedAt     time.Time         `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type GetPayoutRequest struct {
	Search         string `json:"search"`
	CurrencyFilter string `json:"currency"`
	Status         string `json:"status"`
	FromTime       int    `json:"from_time"`
	ToTime         int    `json:"to_time"`
}

func (p *Payout) GetPayoutByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "id = ?", p.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *Payout) GetPayouts(db *gorm.DB, paginator postgresql.Pagination, search string, from int, to int) ([]Payout, postgresql.PaginationResponse, error) {
	details := []Payout{}
	query := ""

	if search != "" {
		query = addQuery(query, fmt.Sprintf("reference = '%v'", search), "and")
	}

	if p.CountryID != 0 {
		query = addQuery(query, fmt.Sprintf("country_id = %v", p.CountryID), "and")
	}

	if p.Status != "" {
		query = addQuery(query, fmt.Sprintf("status = '%v'", p.Status), "and")
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

func (p *Payout) CreatePayout(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payout creation failed: %v", err.Error())
	}
	return nil
}

func (p *Payout) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &p)
	return err
}
