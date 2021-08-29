package inmem

import (
	"github.com/chrismeh/lefty/pkg/products"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

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
