package products

import (
	"time"
)

type Store interface {
	FindAll(filter Filter) ([]Product, error)
	Count() int
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

type Filter struct {
	Page            uint
	ProductsPerPage uint
}
