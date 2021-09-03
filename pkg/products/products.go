package products

import (
	"fmt"
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

func (p Product) String() string {
	return fmt.Sprintf("%s %s", p.Manufacturer, p.Model)
}

type Filter struct {
	Search          string
	Page            uint
	ProductsPerPage uint
}
