package inmem

import (
	"fmt"
	"github.com/chrismeh/lefty/pkg/products"
	"sync"
	"time"
)

type ProductStore struct {
	products map[string]products.Product
	mu       *sync.Mutex
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make(map[string]products.Product),
		mu:       &sync.Mutex{},
	}
}

func (p *ProductStore) Upsert(products []products.Product) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for _, product := range products {
		key := buildProductKey(product)
		if _, exists := p.products[key]; !exists {
			product.CreatedAt = now
		}
		product.UpdatedAt = now
		p.products[key] = product
	}

	return nil
}

func buildProductKey(p products.Product) string {
	return fmt.Sprintf("%s-%s-%s", p.Retailer, p.Manufacturer, p.Model)
}
