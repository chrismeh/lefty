package retailer

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestThomann_LoadProducts(t *testing.T) {
	f, err := os.Open(path.Join("testdata", "thomann_basses_six_string.html"))
	assert.NoError(t, err)
	defer f.Close()

	httpClient := stubHttpClient{getFunc: func(url string) (*http.Response, error) {
		return &http.Response{Body: f}, nil
	}}

	tho := Thomann{http: httpClient}
	response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html")
	assert.NoError(t, err)

	products := response.Products

	assert.Equal(t, uint(1), response.CurrentPage)
	assert.Equal(t, uint(1), response.LastPage)

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

}

type stubHttpClient struct {
	getFunc func(url string) (*http.Response, error)
}

func (s stubHttpClient) Get(url string) (*http.Response, error) {
	return s.getFunc(url)
}
