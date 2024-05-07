package entity

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          string       `db:"id" json:"productId"`
	Name        string       `json:"name" validate:"required,min=1,max=30"`
	SKU         string       `json:"sku" validate:"required,min=1,max=30"`
	Category    string       `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageURL    string       `json:"imageUrl" validate:"required,url"`
	Notes       string       `json:"notes" validate:"required,min=1,max=200"`
	Price       float64      `json:"price" validate:"required,min=1"`
	Stock       int          `json:"stock" validate:"required,min=0,max=100000"`
	Location    string       `json:"location" validate:"required,min=1,max=200"`
	IsAvailable bool         `json:"isAvailable" validate:"required"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

type ProductResponse struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
}
