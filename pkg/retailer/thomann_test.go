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
	t.Parallel()

	t.Run("parse all products on a product page", func(t *testing.T) {
		t.Parallel()

		tho := Thomann{newTestHTTPClientForFixture("thomann_basses_six_strings.html")}
		response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", RequestOptions{})
		assert.NoError(t, err)

		prds := response.Products

		assert.Len(t, prds, 2)
		assert.Equal(t, "Thomann", prds[0].Retailer)
		assert.Equal(t, "ESP", prds[0].Manufacturer)
		assert.Equal(t, "LTD B206SM Natural Satin Left", prds[0].Model)
		assert.Equal(t, "6 saitige Linkshänder E-Bässe", prds[0].Category)
		assert.Equal(t, true, prds[0].IsAvailable)
		assert.Equal(t, "In 4–5 Wochen lieferbar", prds[0].AvailabilityInfo)
		assert.Equal(t, AvailabilityWithinWeeks, prds[0].AvailabilityScore)
		assert.Equal(t, float64(599), prds[0].Price)
		assert.Equal(t, "https://www.thomann.de/de/esp_ltd_b206sm_natural_satin_left_443915.htm?listPosition=0", prds[0].ProductURL)
		assert.Equal(t, "https://thumbs.static-thomann.de/thumb/thumb220x220/pics/prod/443915.jpg", prds[0].ThumbnailURL)

		assert.Equal(t, "Thomann", prds[1].Retailer)
		assert.Equal(t, "Warwick", prds[1].Manufacturer)
		assert.Equal(t, "RB Corvette Basic 6 SBHP LH", prds[1].Model)
		assert.Equal(t, "6 saitige Linkshänder E-Bässe", prds[1].Category)
		assert.Equal(t, true, prds[1].IsAvailable)
		assert.Equal(t, "In 8–10 Wochen lieferbar", prds[1].AvailabilityInfo)
		assert.Equal(t, AvailabilityWithinWeeks, prds[1].AvailabilityScore)
		assert.Equal(t, float64(925), prds[1].Price)
		assert.Equal(t, "https://www.thomann.de/de/warwick_rb_corvette_basic_6_sbhp_lh.htm?listPosition=1", prds[1].ProductURL)
		assert.Equal(t, "https://thumbs.static-thomann.de/thumb/thumb220x220/pics/prod/450435.jpg", prds[1].ThumbnailURL)
	})

	t.Run("parse pagination when there is only a single page", func(t *testing.T) {
		t.Parallel()

		tho := Thomann{newTestHTTPClientForFixture("thomann_basses_six_strings.html")}
		response, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", RequestOptions{})
		assert.NoError(t, err)

		assert.Equal(t, uint(1), response.CurrentPage)
		assert.Equal(t, uint(1), response.LastPage)
	})

	t.Run("parse pagination when there are multiple pages", func(t *testing.T) {
		t.Parallel()

		tho := Thomann{newTestHTTPClientForFixture("thomann_basses_four_strings_second_page.html")}
		response, err := tho.LoadProducts("4_saitige_linkshaender_e-baesse.html", RequestOptions{})
		assert.NoError(t, err)

		assert.Equal(t, uint(2), response.CurrentPage)
		assert.Equal(t, uint(5), response.LastPage)
	})

	t.Run("use correct pagination query parameters depending on RequestOptions struct", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			Name              string
			Page              uint
			ExpectedURLSuffix string
		}{
			{
				Name:              "zero-value RequestOptions",
				Page:              0,
				ExpectedURLSuffix: "?ls=100&pg=1",
			},
			{
				Name:              "valid RequestOptions",
				Page:              2,
				ExpectedURLSuffix: "?ls=100&pg=2",
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				t.Parallel()

				httpSpy := testHTTPClient{
					getFunc: func(url string) (*http.Response, error) {
						return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(""))}, nil
					},
				}
				tho := Thomann{http: &httpSpy}

				options := RequestOptions{Page: tt.Page}
				_, _ = tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", options)

				assert.Equal(t, tt.ExpectedURLSuffix, httpSpy.lastURL[strings.LastIndex(httpSpy.lastURL, "?"):])
			})
		}
	})

	t.Run("return error when page is out of bounds", func(t *testing.T) {
		t.Parallel()

		tho := Thomann{newTestHTTPClientForFixture("thomann_basses_six_strings.html")}

		options := RequestOptions{Page: 1337}
		_, err := tho.LoadProducts("6_saitige_linkshaender_e-baesse.html", options)

		assert.Error(t, err)
	})
}

func TestAvailability_Score(t *testing.T) {
	tests := []struct {
		Name          string
		Status        int
		ExpectedScore int
	}{
		{Name: "product available", Status: 1, ExpectedScore: AvailabilityAvailable},
		{Name: "product available in days", Status: 2, ExpectedScore: AvailabilityWithinDays},
		{Name: "product available in weeks", Status: 4, ExpectedScore: AvailabilityWithinWeeks},
		{Name: "product availability unknown", Status: 1337, ExpectedScore: AvailabilityUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			av := availability{Status: tt.Status}
			assert.Equal(t, tt.ExpectedScore, av.Score())
		})
	}
}

type testHTTPClient struct {
	lastURL string
	getFunc func(url string) (*http.Response, error)
}

func (s *testHTTPClient) Get(url string) (*http.Response, error) {
	s.lastURL = url
	return s.getFunc(url)
}

func newTestHTTPClientForFixture(fixture string) *testHTTPClient {
	testdata, err := ioutil.ReadFile(path.Join("testdata", fixture))
	if err != nil {
		panic(err)
	}

	return &testHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewReader(testdata))}, nil
		},
	}
}
