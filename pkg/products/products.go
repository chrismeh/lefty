package products

import (
	"time"
)

type Product struct {
	Retailer         string    `json:"retailer"`
	Manufacturer     string    `json:"manufacturer"`
	Model            string    `json:"model"`
	Category         string    `json:"category"`
	IsAvailable      bool      `json:"is_available"`
	AvailabilityInfo string    `json:"availability_info"`
	Price            float64   `json:"price"`
	ProductURL       string    `json:"product_url"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Filter struct {
	Page            uint
	ProductsPerPage uint
}
