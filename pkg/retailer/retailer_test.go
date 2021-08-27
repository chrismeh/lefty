package retailer

import (
	"github.com/chrismeh/lefty/pkg/products"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadProducts(t *testing.T) {
	retailer := stubRetailer{}
	retailer.CategoriesFunc = func() []string {
		return []string{"guitars"}
	}

	t.Run("return products for single category without pagination", func(t *testing.T) {
		retailer.LoadProductsFunc = func(category string, options RequestOptions) (ProductResponse, error) {
			pr := ProductResponse{
				Products: []products.Product{
					{
						Retailer:     "Test",
						Manufacturer: "Fender",
						Model:        "AM Pro II Jazzmaster LH MN MYS",
					},
				},
				CurrentPage: 1,
				LastPage:    1,
			}

			return pr, nil
		}

		prds, err := LoadProducts(retailer)
		assert.NoError(t, err)

		assert.Len(t, prds, 1)
		assert.Equal(t, "AM Pro II Jazzmaster LH MN MYS", prds[0].Model)
	})

	t.Run("return products for single category with multiple pages", func(t *testing.T) {
		retailer.LoadProductsFunc = func(category string, options RequestOptions) (ProductResponse, error) {
			var p products.Product
			switch options.Page {
			case 2:
				p = products.Product{Manufacturer: "Gretsch", Model: "G2622LH Strml. DC CB Gunmetal"}
				options.Page = 2
			default:
				p = products.Product{Manufacturer: "Fender", Model: "AM Pro II Jazzmaster LH MN MYS"}
				options.Page = 1
			}

			pr := ProductResponse{Products: []products.Product{p}, CurrentPage: options.Page, LastPage: 2}
			return pr, nil
		}

		prds, err := LoadProducts(retailer)
		assert.NoError(t, err)

		assert.Len(t, prds, 2)
		assert.Equal(t, "AM Pro II Jazzmaster LH MN MYS", prds[0].Model)
		assert.Equal(t, "G2622LH Strml. DC CB Gunmetal", prds[1].Model)
	})

	t.Run("return products for multiple categories with single pages", func(t *testing.T) {
		retailer.CategoriesFunc = func() []string { return []string{"basses", "guitars"} }
		retailer.LoadProductsFunc = func(category string, options RequestOptions) (ProductResponse, error) {
			var p products.Product
			switch category {
			case "basses":
				p = products.Product{Manufacturer: "Fender", Model: "AM Pro II P Bass MN MYS SFG LH"}
			default:
				p = products.Product{Manufacturer: "ESP", Model: "LTD TE-200 Maple STBC LH"}
			}

			pr := ProductResponse{Products: []products.Product{p}, CurrentPage: 1, LastPage: 1}
			return pr, nil
		}

		prds, err := LoadProducts(retailer)
		assert.NoError(t, err)

		assert.Len(t, prds, 2)
		assert.Equal(t, "AM Pro II P Bass MN MYS SFG LH", prds[0].Model)
		assert.Equal(t, "LTD TE-200 Maple STBC LH", prds[1].Model)
	})
}

type stubRetailer struct {
	LoadProductsFunc func(string, RequestOptions) (ProductResponse, error)
	CategoriesFunc   func() []string
}

func (s stubRetailer) LoadProducts(category string, options RequestOptions) (ProductResponse, error) {
	return s.LoadProductsFunc(category, options)
}

func (s stubRetailer) Categories() []string {
	return s.CategoriesFunc()
}
