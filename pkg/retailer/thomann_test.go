package retailer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"
)

func TestThomann_LoadProducts(t *testing.T) {
	t.Run("parse all products on a product page", func(t *testing.T) {
		tho := newThomannForFixture(t, "thomann_basses_six_strings.html")
		response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", RequestOptions{})
		assert.NoError(t, err)

		products := response.Products

		assert.Len(t, products, 2)
		assert.Equal(t, "Thomann", products[0].Retailer)
		assert.Equal(t, "ESP", products[0].Manufacturer)
		assert.Equal(t, "LTD B206SM Natural Satin Left", products[0].Model)
		assert.Equal(t, "6 saitige Linkshänder E-Bässe", products[0].Category)
		assert.Equal(t, true, products[0].IsAvailable)
		assert.Equal(t, "In 1-2 Wochen lieferbar", products[0].AvailabilityInfo)
		assert.Equal(t, float64(599), products[0].Price)
		assert.Equal(t, "https://www.thomann.de/de/esp_ltd_b206sm_natural_satin_left_443915.htm?listPosition=0", products[0].ProductURL)
		assert.Equal(t, "https://thumbs.static-thomann.de/thumb/thumb220x220/pics/prod/443915.jpg", products[0].ThumbnailURL)

		assert.Equal(t, "Thomann", products[1].Retailer)
		assert.Equal(t, "Warwick", products[1].Manufacturer)
		assert.Equal(t, "RB Corvette Basic 6 SBHP LH", products[1].Model)
		assert.Equal(t, "6 saitige Linkshänder E-Bässe", products[1].Category)
		assert.Equal(t, true, products[1].IsAvailable)
		assert.Equal(t, "In 7-9 Wochen lieferbar", products[1].AvailabilityInfo)
		assert.Equal(t, float64(925), products[1].Price)
		assert.Equal(t, "https://www.thomann.de/de/warwick_rb_corvette_basic_6_sbhp_lh.htm?listPosition=1", products[1].ProductURL)
		assert.Equal(t, "https://thumbs.static-thomann.de/thumb/thumb220x220/pics/prod/450435.jpg", products[1].ThumbnailURL)
	})

	t.Run("parse pagination when there is only a single page", func(t *testing.T) {
		tho := newThomannForFixture(t, "thomann_basses_six_strings.html")
		response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", RequestOptions{})
		assert.NoError(t, err)

		assert.Equal(t, uint(1), response.CurrentPage)
		assert.Equal(t, uint(1), response.LastPage)
	})

	t.Run("parse pagination when there are multiple pages", func(t *testing.T) {
		tho := newThomannForFixture(t, "thomann_basses_four_strings_second_page.html")
		response, err := tho.LoadProducts("4_saitige_linkshaender_e-baesse.html", RequestOptions{})
		assert.NoError(t, err)

		assert.Equal(t, uint(2), response.CurrentPage)
		assert.Equal(t, uint(5), response.LastPage)
	})

	t.Run("use correct pagination query parameters depending on RequestOptions struct", func(t *testing.T) {
		tests := []struct {
			Name              string
			ProductsPerPage   uint
			Page              uint
			ExpectedURLSuffix string
		}{
			{
				Name:              "zero-value RequestOptions",
				ProductsPerPage:   0,
				Page:              0,
				ExpectedURLSuffix: "?ls=100&pg=1",
			},
			{
				Name:              "valid RequestOptions",
				ProductsPerPage:   50,
				Page:              2,
				ExpectedURLSuffix: "?ls=50&pg=2",
			},
			{
				Name:              "ProductsPerPage not permitted by retailer",
				ProductsPerPage:   10,
				Page:              1,
				ExpectedURLSuffix: "?ls=100&pg=1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				httpSpy := testHTTPClient{
					getFunc: func(url string) (*http.Response, error) {
						return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
					},
				}
				tho := Thomann{http: &httpSpy}

				options := RequestOptions{ProductsPerPage: tt.ProductsPerPage, Page: tt.Page}
				_, _ = tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", options)

				assert.Equal(t, tt.ExpectedURLSuffix, httpSpy.lastURL[strings.LastIndex(httpSpy.lastURL, "?"):])
			})
		}
	})

	t.Run("return error when page is out of bounds", func(t *testing.T) {
		tho := newThomannForFixture(t, "thomann_basses_six_strings.html")

		options := RequestOptions{ProductsPerPage: 100, Page: 1337}
		_, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", options)

		assert.Error(t, err)
	})
}

type testHTTPClient struct {
	lastURL string
	getFunc func(url string) (*http.Response, error)
}

func (s *testHTTPClient) Get(url string) (*http.Response, error) {
	s.lastURL = url
	return s.getFunc(url)
}

func newThomannForFixture(t *testing.T, filename string) Thomann {
	t.Helper()

	testdata, err := ioutil.ReadFile(path.Join("testdata", filename))
	assert.NoError(t, err)

	httpClient := testHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewReader(testdata))}, nil
		},
	}

	return Thomann{http: &httpClient}
}
