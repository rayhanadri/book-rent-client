package model

import "time"

type Transaction struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	TransactionType string    `json:"transaction_type" gorm:"not null"`
	PaymentMethod   string    `json:"payment_method" gorm:"not null"`
	Amount          int       `json:"amount" gorm:"not null"`
	Status          string    `json:"status" gorm:"not null"`
	Description     string    `json:"description"`
	UserID          int       `json:"user_id,omitempty"`
	User            User      `json:"user" gorm:"foreignKey:UserID"`
	InvoiceID       string    `json:"invoice_id"`
	InvoiceURL      string    `json:"invoice_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	RentID          int       `json:"rent_id,omitempty"`
	Rent            Rent      `json:"rent" gorm:"foreignKey:RentID"`
}
