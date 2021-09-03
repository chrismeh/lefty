package inmem

import (
	"github.com/chrismeh/lefty/pkg/products"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestProductStore_Count(t *testing.T) {
	t.Run("return number of products", func(t *testing.T) {
		productMap := map[string]products.Product{
			"foo": {Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819},
			"bar": {Manufacturer: "Epiphone", Model: "SG Standard Alpine White LH", Price: 449},
		}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		count := store.Count(products.Filter{})
		assert.Equal(t, 2, count)
	})

	t.Run("return number of products that match the filter criteria", func(t *testing.T) {
		productMap := map[string]products.Product{
			"foo": {Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819},
			"bar": {Manufacturer: "Epiphone", Model: "SG Standard Alpine White LH", Price: 449},
		}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		count := store.Count(products.Filter{Search: "Fender"})
		assert.Equal(t, 1, count)
	})
}

func TestProductStore_FindAll(t *testing.T) {
	t.Run("return a slice of products", func(t *testing.T) {
		p := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS"}
		store := ProductStore{products: map[string]products.Product{"foo": p}, mu: &sync.Mutex{}}

		prds, err := store.FindAll(products.Filter{})
		assert.NoError(t, err)

		assert.Len(t, prds, 1)
		assert.Equal(t, "Fender", prds[0].Manufacturer)
		assert.Equal(t, "AM Pro II Jazzmaster LH MN MYS", prds[0].Model)
	})

	t.Run("sort by price ascending by default", func(t *testing.T) {
		p1 := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819}
		p2 := products.Product{Manufacturer: "Fender", Model: "SQ CV 60s Jazzmaster LH LRL OW", Price: 394}
		productMap := map[string]products.Product{"foo": p1, "bar": p2}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		prds, err := store.FindAll(products.Filter{})
		assert.NoError(t, err)

		assert.Equal(t, float64(394), prds[0].Price)
		assert.Equal(t, float64(1819), prds[1].Price)
	})

	t.Run("return a paginated slice of products", func(t *testing.T) {
		p1 := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819}
		p2 := products.Product{Manufacturer: "Fender", Model: "SQ CV 60s Jazzmaster LH LRL OW", Price: 394}
		productMap := map[string]products.Product{"foo": p1, "bar": p2}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		filter := products.Filter{Page: 2, ProductsPerPage: 1}
		prds, err := store.FindAll(filter)
		assert.NoError(t, err)

		assert.Len(t, prds, 1)
		assert.Equal(t, float64(1819), prds[0].Price)
	})

	t.Run("apply default pagination settings when pagination data is invalid", func(t *testing.T) {
		productMap := map[string]products.Product{
			"foo": {Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819},
			"bar": {Manufacturer: "Fender", Model: "SQ CV 60s Jazzmaster LH LRL OW", Price: 394},
		}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		tests := []struct {
			Name            string
			Page            uint
			ProductsPerPage uint
		}{
			{Name: "products per page is larger than product count", Page: 1, ProductsPerPage: 50},
			{Name: "page is zero", Page: 0, ProductsPerPage: 50},
			{Name: "page does not exist", Page: 100, ProductsPerPage: 50},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				filter := products.Filter{Page: tt.Page, ProductsPerPage: tt.ProductsPerPage}
				prds, err := store.FindAll(filter)
				assert.NoError(t, err)

				assert.Len(t, prds, 2)
			})
		}
	})

	t.Run("return slice of remaining products on the last page", func(t *testing.T) {
		productMap := map[string]products.Product{
			"foo": {Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819},
			"bar": {Manufacturer: "Fender", Model: "SQ CV 60s Jazzmaster LH LRL OW", Price: 394},
			"baz": {Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH 3TSB", Price: 1799},
		}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		prds, err := store.FindAll(products.Filter{Page: 2, ProductsPerPage: 2})
		assert.NoError(t, err)

		assert.Len(t, prds, 1)
		assert.Equal(t, "AM Pro II Jazzmaster LH MN MYS", prds[0].Model)
	})

	t.Run("return only products that match the filter criteria", func(t *testing.T) {
		productMap := map[string]products.Product{
			"foo": {Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819},
			"bar": {Manufacturer: "Epiphone", Model: "SG Standard Alpine White LH", Price: 449},
		}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		tests := []struct {
			Name          string
			Search        string
			ExpectedModel string
		}{
			{
				Name:          "find product by model",
				Search:        "SG",
				ExpectedModel: "SG Standard Alpine White LH",
			},
			{
				Name:          "find product by manufacturer",
				Search:        "Fender",
				ExpectedModel: "AM Pro II Jazzmaster LH MN MYS",
			},
			{
				Name:          "find product by model, case insensitive",
				Search:        "sg",
				ExpectedModel: "SG Standard Alpine White LH",
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				prds, err := store.FindAll(products.Filter{Search: tt.Search})
				assert.NoError(t, err)

				assert.Len(t, prds, 1)
				assert.Equal(t, tt.ExpectedModel, prds[0].Model)
			})
		}
	})
}

func TestProductStore_Upsert(t *testing.T) {
	t.Run("save a new product", func(t *testing.T) {
		p := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS"}
		store := NewProductStore()

		err := store.Upsert([]products.Product{p})

		assert.NoError(t, err)
		assert.Len(t, store.products, 1)

		pk := buildProductKey(p)
		assert.Equal(t, "Fender", store.products[pk].Manufacturer)
		assert.Equal(t, "AM Pro II Jazzmaster LH MN MYS", store.products[pk].Model)
	})

	t.Run("set timestamps when saving a new product", func(t *testing.T) {
		p := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS"}
		store := NewProductStore()

		_ = store.Upsert([]products.Product{p})

		pk := buildProductKey(p)
		assert.Equal(t, store.products[pk].CreatedAt, store.products[pk].UpdatedAt)
		assert.False(t, store.products[pk].CreatedAt.IsZero())
		assert.False(t, store.products[pk].UpdatedAt.IsZero())
	})

	t.Run("setUpdatedAt timestamp when saving an existing product", func(t *testing.T) {
		p := products.Product{
			Manufacturer: "Fender",
			Model:        "AM Pro II Jazzmaster LH MN MYS",
			CreatedAt:    time.Date(2014, 8, 6, 23, 0, 0, 0, time.UTC),
			UpdatedAt:    time.Date(2014, 8, 6, 23, 0, 0, 0, time.UTC),
		}
		pk := buildProductKey(p)
		store := ProductStore{products: map[string]products.Product{pk: p}, mu: &sync.Mutex{}}

		_ = store.Upsert([]products.Product{p})

		assert.NotEqual(t, store.products[pk].CreatedAt, store.products[pk].UpdatedAt)

	})
}
