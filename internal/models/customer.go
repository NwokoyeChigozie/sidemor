package models

import "time"

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
