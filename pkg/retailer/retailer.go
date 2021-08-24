package retailer

import "github.com/chrismeh/lefty/pkg/products"

type ProductResponse struct {
	Products    []products.Product
	CurrentPage uint
	LastPage    uint
}
