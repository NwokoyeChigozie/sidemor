package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Withdrawal struct {
	ID             uint                 `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	MerchantID     int64                `gorm:"column:merchant_id; type:int" json:"merchant_id"`
	Merchant       external_models.User `gorm:"-" json:"merchant"`
	Currency       string               `gorm:"column:currency; type:varchar(255)" json:"currency"`
	Amount         float64              `gorm:"column:amount; type:decimal(20,2)" json:"amount"`
	WithdrawalDate time.Time            `gorm:"column:withdrawal_date; autoCreateTime" json:"withdrawal_date"`
	Status         TransactionStatus    `gorm:"column:status; type:varchar(255)" json:"status"`
	CreatedAt      time.Time            `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time            `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type RequestWithdrawalRequest struct {
	Currency       string  `json:"currency"  validate:"required"`
	Amount         float64 `json:"amount"  validate:"required"`
	WithdrawalDate int     `json:"withdrawal_date" validate:"required"`
}

type GetWithdrawalRequest struct {
	Search         string `json:"search"`
	CurrencyFilter string `json:"currency"`
	Status         string `json:"status"`
	FromTime       int    `json:"from_time"`
	ToTime         int    `json:"to_time"`
}

func (w *Withdrawal) GetWithdrawalByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &w, "id = ?", w.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (w *Withdrawal) GetWithdrawals(db *gorm.DB, paginator *postgresql.Pagination, userIds []int, from int, to int) ([]Withdrawal, postgresql.PaginationResponse, error) {
	var (
		details    = []Withdrawal{}
		query      = ""
		pagination postgresql.PaginationResponse
	)

	if len(userIds) > 0 {
		idsString := []string{}
		for _, v := range userIds {
			idsString = append(idsString, fmt.Sprintf("%v", v))
		}
		query = addQuery(query, fmt.Sprintf(" merchant_id IN (%s)", strings.Join(idsString, ",")), "and")
	}

	if w.MerchantID != 0 {
		query = addQuery(query, fmt.Sprintf("merchant_id = %v", w.MerchantID), "and")
	}

	if w.Status != "" {
		query = addQuery(query, fmt.Sprintf("status = '%v'", w.Status), "and")
	}

	if from != 0 {
		fromTime := time.Unix(int64(from), 0)
		query = addQuery(query, fmt.Sprintf("withdrawal_date >= '%v'", fromTime), "and")
	}

	if to != 0 {
		toTime := time.Unix(int64(to), 0)
		query = addQuery(query, fmt.Sprintf("withdrawal_date >= '%v'", toTime), "and")
	}

	if paginator == nil {
		err := postgresql.SelectAllFromDbOrderBy(db, "withdrawal_date", "asc", &details, query)
		if err != nil {
			return details, pagination, err
		}
	} else {
		pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "withdrawal_date", "asc", *paginator, &details, query)
		if err != nil {
			return details, pagination, err
		}
	}

	return details, pagination, nil
}

func (w *Withdrawal) CreateWithdrawal(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("Withdrawal creation failed: %v", err.Error())
	}
	return nil
}

func (w *Withdrawal) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &w)
	return err
}
