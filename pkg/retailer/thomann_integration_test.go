//go:build integration
// +build integration

package retailer

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestThomann_LoadProductsIntegration(t *testing.T) {
	t.Parallel()

	c := &http.Client{Timeout: 5 * time.Second}
	tho := Thomann{http: c}

	response, err := tho.LoadProducts(tho.Categories()[0], RequestOptions{})
	assert.NoError(t, err)

	assert.Len(t, response.Products, 100)

	p := response.Products[0]
	assert.Equal(t, "Thomann", p.Retailer)
	assert.NotZero(t, p.Manufacturer)
	assert.NotZero(t, p.Model)
	assert.NotZero(t, p.Category)
	assert.NotZero(t, p.AvailabilityInfo)
	assert.NotZero(t, p.Price)
	assert.NotZero(t, p.ProductURL)
	assert.NotZero(t, p.ThumbnailURL)

}
