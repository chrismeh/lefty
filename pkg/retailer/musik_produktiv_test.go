package retailer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMusikProduktiv_LoadProducts(t *testing.T) {
	t.Run("parse product from product page", func(t *testing.T) {
		mp := MusikProduktiv{http: newTestHTTPClientForFixture("musikproduktiv_guitars_eight_strings.html")}

		response, err := mp.LoadProducts("e-gitarre-linkshaender")
		assert.NoError(t, err)

		assert.Len(t, response.Products, 1)
		assert.Equal(t, "Musik Produktiv", response.Products[0].Retailer)
		assert.Equal(t, "Schecter", response.Products[0].Manufacturer)
		assert.Equal(t, "C-8 Deluxe LH SBK", response.Products[0].Model)
		assert.Equal(t, "E-Gitarre (Linkshänder), 8-saitig", response.Products[0].Category)
		assert.Equal(t, false, response.Products[0].IsAvailable)
		assert.Equal(t, "", response.Products[0].AvailabilityInfo)
		assert.Equal(t, float64(599), response.Products[0].Price)
		assert.Equal(t, "https://www.musik-produktiv.de/schecter-c-8-deluxe-lh-sbk.html", response.Products[0].ProductURL)
		assert.Equal(t, "https://sc1.musik-produktiv.com/pic-010125643l/schecter-c-8-deluxe-lh-sbk.jpg", response.Products[0].ThumbnailURL)
	})

	t.Run("parse model and manufacturer titles when manufacturer name contains spaces", func(t *testing.T) {
		mp := MusikProduktiv{http: newTestHTTPClientForFixture("musikproduktiv_guitars_second_page.html")}

		response, err := mp.LoadProducts("e-gitarre-linkshaender")
		assert.NoError(t, err)

		assert.Len(t, response.Products, 20)
		assert.Equal(t, "ESP LTD", response.Products[0].Manufacturer)
		assert.Equal(t, "Signature Iron Cross J.Hetfield Lefthand", response.Products[0].Model)
	})

	t.Run("parse pagination when there is only a single page", func(t *testing.T) {
		mp := MusikProduktiv{http: newTestHTTPClientForFixture("musikproduktiv_guitars_eight_strings.html")}

		response, err := mp.LoadProducts("e-gitarre-linkshaender")
		assert.NoError(t, err)

		assert.Equal(t, uint(1), response.CurrentPage)
		assert.Equal(t, uint(1), response.LastPage)
	})

	t.Run("parse pagination when there are multiple pages", func(t *testing.T) {
		mp := MusikProduktiv{http: newTestHTTPClientForFixture("musikproduktiv_guitars_second_page.html")}

		response, err := mp.LoadProducts("e-gitarre-linkshaender")
		assert.NoError(t, err)

		assert.Equal(t, uint(2), response.CurrentPage)
		assert.Equal(t, uint(6), response.LastPage)
	})

	t.Run("parse pagination on the last page", func(t *testing.T) {
		mp := MusikProduktiv{http: newTestHTTPClientForFixture("musikproduktiv_guitars_last_page.html")}

		response, err := mp.LoadProducts("e-gitarre-linkshaender")
		assert.NoError(t, err)

		assert.Equal(t, uint(6), response.CurrentPage)
		assert.Equal(t, uint(6), response.LastPage)
	})
}
