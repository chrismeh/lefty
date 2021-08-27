package retailer

import "github.com/chrismeh/lefty/pkg/products"

type Retailer interface {
	LoadProducts(category string, options RequestOptions) (ProductResponse, error)
}

func LoadProducts(r Retailer) ([]products.Product, error) {
	var page uint = 1
	resp, err := r.LoadProducts("guitars", RequestOptions{Page: page})
	if err != nil {
		return nil, err
	}

	prds := make([]products.Product, len(resp.Products))
	copy(prds, resp.Products)

	for page < resp.LastPage {
		page++
		resp, err = r.LoadProducts("guitars", RequestOptions{Page: page})
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
