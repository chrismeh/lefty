package retailer

import (
	"encoding/json"
	"fmt"
	"github.com/chrismeh/lefty/pkg/products"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Thomann struct {
	http httpGetter
}

func NewThomann(http httpGetter) Thomann {
	return Thomann{http: http}
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

	re := regexp.MustCompile(`(?ms)({"headline":.+?"})\);`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return ProductResponse{}, fmt.Errorf("unexpected response body structure")
	}

	var p page
	err = json.NewDecoder(strings.NewReader(match[1])).Decode(&p)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	if uint(p.Pagination.LastPage) < options.Page {
		return ProductResponse{}, fmt.Errorf("page %d out of bounds, last page is %d", options.Page, p.Pagination.LastPage)
	}

	productResponse := ProductResponse{
		Products:    p.products(),
		CurrentPage: uint(p.Pagination.CurrentPage),
		LastPage:    uint(p.Pagination.LastPage),
	}
	return productResponse, nil
}

func (t Thomann) Categories() []string {
	return []string{
		"linkshaender_modelle.html",
		"linkshaender_konzertgitarren.html",
		"linkshaender_akustikgitarren.html",
		"4_saitige_linkshaender_e-baesse.html",
		"5_saitige_linkshaender_e-baesse.html",
		"6_saitige_linkshaender_e-baesse.html",
	}
}

func (t Thomann) buildURL(category string, options RequestOptions) string {
	var productsPerPage uint = 100
	var page uint = 1

	if options.Page > 0 {
		page = options.Page
	}

	return fmt.Sprintf("https://www.thomann.de/de/%s?ls=%d&pg=%d", category, productsPerPage, page)
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

		productURL := fmt.Sprintf("https://www.thomann.de/de/%s", v.Link)
		thumbnailURL := fmt.Sprintf("https://thumbs.static-thomann.de/thumb/thumb220x220/pics/prod/%s", v.Image.Name)

		pr[k] = products.Product{
			Retailer:          "Thomann",
			Manufacturer:      v.Manufacturer,
			Model:             v.Model,
			Category:          p.Title,
			IsAvailable:       v.Availability.IsAvailable,
			AvailabilityInfo:  v.Availability.Text,
			AvailabilityScore: v.Availability.Score(),
			Price:             price,
			ProductURL:        productURL,
			ThumbnailURL:      thumbnailURL,
		}
	}

	return pr
}

type articleList struct {
	Articles []article `json:"articles"`
}

type article struct {
	Manufacturer string       `json:"manufacturer"`
	Model        string       `json:"model"`
	Availability availability `json:"availability"`
	Price        price        `json:"price"`
	Link         string       `json:"relativeLink"`
	Image        image        `json:"mainImage"`
}

type availability struct {
	Status      int    `json:"code"`
	IsAvailable bool   `json:"isAvailable"`
	Text        string `json:"textShort"`
}

func (a availability) Score() int {
	switch a.Status {
	case 1:
		return products.AvailabilityAvailable
	case 2:
		return products.AvailabilityWithinDays
	case 4:
		return products.AvailabilityWithinWeeks
	default:
		return products.AvailabilityUnknown
	}
}

type price struct {
	Primary priceEntry `json:"primary"`
}

type priceEntry struct {
	Raw string `json:"rawPrice"`
}

type image struct {
	Name string `json:"fileName"`
}

type pagination struct {
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
}
