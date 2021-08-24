package retailer

import (
	"encoding/json"
	"fmt"
	"github.com/chrismeh/lefty/pkg/products"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Thomann struct {
	http interface {
		Get(url string) (*http.Response, error)
	}
	baseURL string
}

func (t Thomann) LoadProducts(category string) ([]products.Product, error) {
	resp, err := t.http.Get(category)
	if err != nil {
		return nil, fmt.Errorf("could not fetch products from thomann.de: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body from thomann.de: %w", err)
	}

	re := regexp.MustCompile(`(?ms)({"headline":.+?"})\], {"general`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return nil, fmt.Errorf("unexpected response body structure")
	}

	var p page
	err = json.NewDecoder(strings.NewReader(match[1])).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return p.products(), nil
}

type page struct {
	Title       string      `json:"headline"`
	ArticleList articleList `json:"articleListsSettings"`
}

func (p page) products() []products.Product {
	pr := make([]products.Product, len(p.ArticleList.Articles))
	for k, v := range p.ArticleList.Articles {
		price, err := strconv.ParseFloat(v.Price.Primary.Raw, 32)
		if err != nil {
			price = 0.00
		}

		var thumbnailURL string
		if v.Image.Exists {
			thumbnailURL = fmt.Sprintf("https://thumbs.static-thomann.de/thumb/thumb220x220/pics/prod/%s", v.Image.Name)
		}

		pr[k] = products.Product{
			Retailer:         "Thomann",
			Manufacturer:     v.Manufacturer,
			Model:            v.Model,
			Category:         p.Title,
			IsAvailable:      v.Availability.IsAvailable,
			AvailabilityInfo: v.Availability.Text,
			Price:            price,
			ProductURL:       v.Link,
			ThumbnailURL:     thumbnailURL,
		}
	}

	return pr
}

type articleList struct {
	Articles []article `json:"articles"`
}

type article struct {
	Manufacturer string       `json:"manufacturer"`
	Model        string       `json:"name"`
	Availability availability `json:"availability"`
	Price        price        `json:"price"`
	Link         string       `json:"link"`
	Image        image        `json:"image"`
}

type availability struct {
	IsAvailable bool   `json:"isAvailable"`
	Text        string `json:"text"`
}

type price struct {
	Primary priceEntry `json:"primary"`
}

type priceEntry struct {
	Raw string `json:"raw"`
}

type image struct {
	Name   string `json:"fname"`
	Exists bool   `json:"exists"`
}
