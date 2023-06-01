package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Rank      int32     `json:"rank"`
}
