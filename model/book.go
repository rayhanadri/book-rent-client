package model

import "time"

type Book struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Author      string    `json:"author" gorm:"not null"`
	Publisher   string    `json:"publisher" gorm:"not null"`
	PublishedAt time.Time `json:"published_at" gorm:"not null"`
	ISBN        string    `json:"isbn" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	Available   bool      `json:"available" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	Description string    `json:"description"`
}
