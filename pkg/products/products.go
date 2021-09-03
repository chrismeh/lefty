package products

import (
	"time"
)

type Store interface {
	FindAll() ([]Product, error)
	Upsert(products []Product) error
}

type Product struct {
	Retailer         string
	Manufacturer     string
	Model            string
	Category         string
	IsAvailable      bool
	AvailabilityInfo string
	Price            float64
	ProductURL       string
	ThumbnailURL     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
