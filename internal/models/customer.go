package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Customer struct {
	ID                uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID         int64     `gorm:"column:account_id; type:int; not null; comment: from the API key used to authenticate the call from merchant" json:"account_id"`
	Email             string    `gorm:"column:email; type:varchar(255)" json:"email"`
	Firstname         string    `gorm:"column:firstname; type:varchar(255)" json:"firstname"`
	Lastname          string    `gorm:"column:lastname; type:varchar(255)" json:"lastname"`
	Address           string    `gorm:"column:address; type:varchar(255)" json:"address"`
	City              string    `gorm:"column:city; type:varchar(255)" json:"city"`
	State             string    `gorm:"column:state; type:varchar(255)" json:"state"`
	PhoneNumber       string    `gorm:"column:phone_number; type:varchar(255)" json:"phone_number"`
	CountryID         int64     `gorm:"column:country_id; type:int" json:"country_id"`
	NumberOfPayments  int64     `gorm:"column:number_of_payments; type:int;default:0" json:"number_of_payments"`
	LastPaymentMadeAt time.Time `gorm:"column:last_payment_made_at" json:"last_payment_made_at"`
	CreatedAt         time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (c *Customer) CreateCustomer(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &c)
	if err != nil {
		return fmt.Errorf("customer creation failed: %v", err.Error())
	}
	return nil
}

func (c *Customer) GetCustomerByAccountIDAndEmail(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &c, "account_id = ? and email=?", c.AccountID, c.Email)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (c *Customer) GetCustomerByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &c, "id = ?", c.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (c *Customer) GetCustomers(db *gorm.DB, paginator postgresql.Pagination, search string) ([]Customer, postgresql.PaginationResponse, error) {
	var (
		details    = []Customer{}
		extraQuery string
	)

	if search != "" {
		extraQuery += fmt.Sprintf(` and (
			Lower(email) = '%[1]v'
			or Lower(firstname) = '%[1]v'
			or Lower(lastname) = '%[1]v'
			or Lower(address) = '%[1]v'
			or Lower(city) = '%[1]v'
			or Lower(state) = '%[1]v'
			or Lower(phone_number) = '%[1]v'
			or Lower(country_id) = '%[1]v'
			or Lower(number_of_payments) = '%[1]v'
			)`, strings.ToLower(search))
	}
	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, fmt.Sprintf("account_id = ? %v", extraQuery), c.AccountID)
	if err != nil {
		return details, pagination, err
	}
	return details, pagination, nil
}

func (c *Customer) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &c)
	return err
}
