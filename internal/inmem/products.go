package inmem

import (
	"encoding/json"
	"fmt"
	"github.com/chrismeh/lefty/pkg/products"
	"io"
	"math"
	"sort"
	"strings"
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

func (p *ProductStore) FindAll(f products.Filter) ([]products.Product, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	prds := make([]products.Product, 0, len(p.products))
	for _, v := range p.products {
		if f.Search == "" || productMatchesFilter(v, f) {
			prds = append(prds, v)
		}
	}
	if len(prds) == 0 {
		return prds, nil
	}

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].Price < prds[j].Price
	})

	return paginate(prds, f), nil
}

func (p *ProductStore) Count(f products.Filter) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	var count int
	for _, v := range p.products {
		if f.Search == "" || productMatchesFilter(v, f) {
			count++
		}
	}

	return count
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

func (p *ProductStore) Dump(w io.Writer) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return json.NewEncoder(w).Encode(p.products)
}

func (p *ProductStore) Load(r io.Reader) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return json.NewDecoder(r).Decode(&p.products)
}

func productMatchesFilter(p products.Product, f products.Filter) bool {
	name := strings.ToLower(p.String())
	search := strings.ToLower(f.Search)

	return strings.Contains(name, search)
}

func paginate(prds []products.Product, f products.Filter) []products.Product {
	count := uint(len(prds))
	if f.Page == 0 {
		f.Page = 1
	}
	if f.ProductsPerPage == 0 {
		f.ProductsPerPage = 50
	}
	if f.ProductsPerPage > count {
		f.ProductsPerPage = count
	}

	lastPage := math.Ceil(float64(count) / float64(f.ProductsPerPage))
	if f.Page > uint(lastPage) {
		f.Page = 1
	}

	offset := (f.Page - 1) * f.ProductsPerPage
	limit := f.ProductsPerPage

	if offset+limit > count {
		limit = count - offset
	}

	return prds[offset : offset+limit]
}

func buildProductKey(p products.Product) string {
	return fmt.Sprintf("%s-%s-%s", p.Retailer, p.Manufacturer, p.Model)
}
