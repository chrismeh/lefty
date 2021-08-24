package retailer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"path"
	"testing"
)

func TestThomann_LoadProducts(t *testing.T) {
	t.Run("parse all products on a product page", func(t *testing.T) {
		tho := newThomannForFixture(t, "thomann_basses_six_strings.html")
		response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html")
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
		response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html")
		assert.NoError(t, err)

		assert.Equal(t, uint(1), response.CurrentPage)
		assert.Equal(t, uint(1), response.LastPage)
	})

	t.Run("parse pagination when there are multiple pages", func(t *testing.T) {
		tho := newThomannForFixture(t, "thomann_basses_four_strings_second_page.html")
		response, err := tho.LoadProducts("4_saitige_linkshaender_e-baesse.html")
		assert.NoError(t, err)

		assert.Equal(t, uint(2), response.CurrentPage)
		assert.Equal(t, uint(5), response.LastPage)
	})

	t.Run("use correct pagination settings when making the HTTP request", func(t *testing.T) {
		httpSpy := testHTTPClient{
			requestedURLs: make([]string, 0),
			getFunc: func(url string) (*http.Response, error) {
				return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
			},
		}
		tho := Thomann{
			http:    &httpSpy,
			baseURL: "https://www.thomann.de/de",
		}

		_, _ = tho.LoadProducts("6_saitige_linkshaender_e-baesse.html")

		assert.Len(t, httpSpy.requestedURLs, 1)
		assert.Equal(t, "https://www.thomann.de/de/6_saitige_linkshaender_e-baesse.html?ls=100&pg=1", httpSpy.requestedURLs[0])
	})
}

type testHTTPClient struct {
	requestedURLs []string
	getFunc       func(url string) (*http.Response, error)
}

func (s *testHTTPClient) Get(url string) (*http.Response, error) {
	s.requestedURLs = append(s.requestedURLs, url)
	return s.getFunc(url)
}

func newThomannForFixture(t *testing.T, filename string) Thomann {
	t.Helper()

	testdata, err := ioutil.ReadFile(path.Join("testdata", filename))
	assert.NoError(t, err)

	httpClient := testHTTPClient{
		requestedURLs: make([]string, 0),
		getFunc: func(url string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewReader(testdata))}, nil
		},
	}

	return Thomann{http: &httpClient}
}
