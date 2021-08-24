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
	products, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html")
	assert.NoError(t, err)

	assert.Len(t, products, 2)
	assert.Equal(t, "ESP", products[0].Manufacturer)
	assert.Equal(t, "LTD B206SM Natural Satin Left", products[0].Model)
	assert.Equal(t, "Warwick", products[1].Manufacturer)
	assert.Equal(t, "RB Corvette Basic 6 SBHP LH", products[1].Model)
}

type stubHttpClient struct {
	getFunc func(url string) (*http.Response, error)
}

func (s stubHttpClient) Get(url string) (*http.Response, error) {
	return s.getFunc(url)
}
