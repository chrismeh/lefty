package retailer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestMusikProduktiv_LoadProducts(t *testing.T) {
	testdata, err := os.ReadFile(path.Join("testdata", "musikproduktiv_guitar_eight_strings.html"))
	assert.NoError(t, err)

	httpStub := testHTTPClient{getFunc: func(url string) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewReader(testdata))}, nil
	}}
	mp := MusikProduktiv{&httpStub}

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
}
