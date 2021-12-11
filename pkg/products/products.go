package products

import (
	"fmt"
	"time"
)

const (
	OrderPriceDesc          string = "-price"
	OrderByAvailabilityAsc         = "availability"
	OrderByAvailabilityDesc        = "-availability"
	AvailabilityAvailable   int    = 1
	AvailabilityWithinDays         = 2
	AvailabilityWithinWeeks        = 3
	AvailabilityUnknown            = 4
)

type Product struct {
	Retailer          string    `json:"retailer"`
	Manufacturer      string    `json:"manufacturer"`
	Model             string    `json:"model"`
	Category          string    `json:"category"`
	IsAvailable       bool      `json:"is_available"`
	AvailabilityInfo  string    `json:"availability_info"`
	AvailabilityScore int       `json:"availability_score"`
	Price             float64   `json:"price"`
	ProductURL        string    `json:"product_url"`
	ThumbnailURL      string    `json:"thumbnail_url"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (p Product) String() string {
	return fmt.Sprintf("%s %s", p.Manufacturer, p.Model)
}

type Filter struct {
	Search          string
	OrderBy         string
	Retailer        string
	Page            uint
	ProductsPerPage uint
}

func (f Filter) HasFilterCriteria() bool {
	return f.Search != "" || f.Retailer != ""
}
