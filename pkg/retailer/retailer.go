package retailer

import "github.com/chrismeh/lefty/pkg/products"

type Retailer interface {
	LoadProducts(category string, options RequestOptions) (ProductResponse, error)
	Categories() []string
}

func LoadProducts(r Retailer) ([]products.Product, error) {
	prds := make([]products.Product, 0)

	for _, category := range r.Categories() {
		categoryProducts, err := loadProductsFromCategory(r, category)
		if err != nil {
			return nil, err
		}

		prds = append(prds, categoryProducts...)
	}

	return prds, nil
}

func loadProductsFromCategory(r Retailer, category string) ([]products.Product, error) {
	var page uint = 1
	resp, err := r.LoadProducts(category, RequestOptions{Page: page})
	if err != nil {
		return nil, err
	}

	prds := make([]products.Product, len(resp.Products))
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
	Products    []products.Product
	CurrentPage uint
	LastPage    uint
}
