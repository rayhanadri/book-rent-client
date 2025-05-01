package model

import (
	"time"
)

type Rent struct {
	ID            int       `json:"id,omitempty" gorm:"primaryKey"`
	BookID        int       `json:"book_id" gorm:"not null"`
	UserID        int       `json:"user_id" gorm:"not null"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	TotalPrice    int       `json:"total_price" gorm:"not null"`
	RentStartDate time.Time `json:"rent_start_date"`
	RentEndDate   time.Time `json:"rent_end_date"`
	RentStatus    string    `json:"rent_status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
