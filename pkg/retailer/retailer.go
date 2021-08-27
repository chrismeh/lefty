package retailer

import "github.com/chrismeh/lefty/pkg/products"

type Retailer interface {
	LoadProducts(category string, options RequestOptions) (ProductResponse, error)
}

func LoadProducts(r Retailer) ([]products.Product, error) {
	resp, err := r.LoadProducts("guitars", RequestOptions{})
	if err != nil {
		return nil, err
	}

	return resp.Products, nil
}

type RequestOptions struct {
	Page uint
}

type ProductResponse struct {
	Products    []products.Product
	CurrentPage uint
	LastPage    uint
}
