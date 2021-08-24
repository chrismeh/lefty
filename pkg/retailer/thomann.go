package retailer

import (
	"encoding/json"
	"fmt"
	"github.com/chrismeh/lefty/pkg/products"
	"io/ioutil"
	"net/http"
	"regexp"
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
	ArticleList articleList `json:"articleListsSettings"`
}

func (p page) products() []products.Product {
	pr := make([]products.Product, len(p.ArticleList.Articles))
	for k, v := range p.ArticleList.Articles {
		pr[k] = products.Product{
			Manufacturer: v.Manufacturer,
			Model:        v.Model,
		}
	}

	return pr
}

type articleList struct {
	Articles []article `json:"articles"`
}

type article struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"name"`
}
