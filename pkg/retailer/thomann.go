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

func (t Thomann) LoadProducts(category string, options RequestOptions) (ProductResponse, error) {
	resp, err := t.http.Get(t.buildURL(category, options))
	if err != nil {
		return ProductResponse{}, fmt.Errorf("could not fetch products from thomann.de: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("could not read response body from thomann.de: %w", err)
	}

	re := regexp.MustCompile(`(?ms)({"headline":.+?"})\], {"general`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return ProductResponse{}, fmt.Errorf("unexpected response body structure")
	}

	var p page
	err = json.NewDecoder(strings.NewReader(match[1])).Decode(&p)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	productResponse := ProductResponse{
		Products:    p.products(),
		CurrentPage: uint(p.Pagination.CurrentPage),
		LastPage:    uint(p.Pagination.LastPage),
	}
	return productResponse, nil
}

func (t Thomann) buildURL(category string, options RequestOptions) string {
	var productsPerPage uint = 100
	var page uint = 1

	validProductsPerPage := []uint{25, 50, 100}
	for _, v := range validProductsPerPage {
		if options.ProductsPerPage == v {
			productsPerPage = v
			break
		}
	}

	if options.Page > 0 {
		page = options.Page
	}

	return fmt.Sprintf("%s/%s?ls=%d&pg=%d", t.baseURL, category, productsPerPage, page)
}

type page struct {
	Title       string      `json:"headline"`
	ArticleList articleList `json:"articleListsSettings"`
	Pagination  pagination  `json:"pagingSettings"`
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

type pagination struct {
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
}
