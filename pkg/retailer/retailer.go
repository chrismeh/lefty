package retailer

import "github.com/chrismeh/lefty/pkg/products"

type RequestOptions struct {
	ProductsPerPage uint
	Page            uint
}

type ProductResponse struct {
	Products    []products.Product
	CurrentPage uint
	LastPage    uint
}
