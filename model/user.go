package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password,omitempty" validate:"required"`
	Role      string    `json:"role,omitempty"`
	Status    string    `json:"status"`
	Balance   int       `json:"balance,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
