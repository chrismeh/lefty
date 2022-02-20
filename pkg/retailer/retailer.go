package retailer

import (
	"net/http"
)

type Retailer interface {
	LoadProducts(category string, options RequestOptions) (ProductResponse, error)
	Categories() []string
}

type ProductUpserter interface {
	Upsert([]Product) error
}

func UpdateRetailers(ps ProductUpserter, retailer ...Retailer) error {
	prds := make([]Product, 0)

	for _, r := range retailer {
		p, err := LoadProducts(r)
		if err != nil {
			return err
		}

		prds = append(prds, p...)
	}

	return ps.Upsert(prds)
}

func LoadProducts(r Retailer) ([]Product, error) {
	prds := make([]Product, 0)

	for _, category := range r.Categories() {
		categoryProducts, err := loadProductsFromCategory(r, category)
		if err != nil {
			return nil, err
		}

		prds = append(prds, categoryProducts...)
	}

	return prds, nil
}

func loadProductsFromCategory(r Retailer, category string) ([]Product, error) {
	var page uint = 1
	resp, err := r.LoadProducts(category, RequestOptions{Page: page})
	if err != nil {
		return nil, err
	}

	prds := make([]Product, len(resp.Products))
	copy(prds, resp.Products)

	for page < resp.LastPage {
		page++
		resp, err = r.LoadProducts(category, RequestOptions{Page: page})
		if err != nil {
			return nil, err
		}

		prds = append(prds, resp.Products...)
	}

	return prds, nil
}

type RequestOptions struct {
	Page uint
}

type ProductResponse struct {
	Products    []Product
	CurrentPage uint
	LastPage    uint
}

type httpGetter interface {
	Get(url string) (*http.Response, error)
}
