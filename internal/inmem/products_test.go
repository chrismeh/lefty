package inmem

import (
	"github.com/chrismeh/lefty/pkg/products"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestProductStore_FindAll(t *testing.T) {
	t.Run("return a slice of products", func(t *testing.T) {
		p := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS"}
		store := ProductStore{products: map[string]products.Product{buildProductKey(p): p}, mu: &sync.Mutex{}}

		prds, err := store.FindAll()
		assert.NoError(t, err)

		assert.Len(t, prds, 1)
		assert.Equal(t, "Fender", prds[0].Manufacturer)
		assert.Equal(t, "AM Pro II Jazzmaster LH MN MYS", prds[0].Model)
	})

	t.Run("should sort by price ascending by default", func(t *testing.T) {
		p1 := products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS", Price: 1819}
		p2 := products.Product{Manufacturer: "Fender", Model: "SQ CV 60s Jazzmaster LH LRL OW", Price: 394}
		productMap := map[string]products.Product{buildProductKey(p1): p1, buildProductKey(p2): p2}
		store := ProductStore{products: productMap, mu: &sync.Mutex{}}

		prds, err := store.FindAll()
		assert.NoError(t, err)

		assert.Equal(t, float64(394), prds[0].Price)
		assert.Equal(t, float64(1819), prds[1].Price)
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
