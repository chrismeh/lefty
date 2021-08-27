package retailer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chrismeh/lefty/pkg/products"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type MusikProduktiv struct {
	http interface {
		Get(url string) (*http.Response, error)
	}
}

func (m MusikProduktiv) LoadProducts(category string) (ProductResponse, error) {
	resp, err := m.http.Get(fmt.Sprintf("https://www.musik-produktiv.de/%s", category))
	if err != nil {
		return ProductResponse{}, fmt.Errorf("could not fetch products from musik-produktiv.de: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("could not create goquery document from reader: %w", err)
	}

	categoryName := doc.Find("div.list_title h1").Text()

	nodes := doc.Find("ul.artgrid li")
	instruments := make([]products.Product, len(nodes.Nodes))

	nodes.Each(func(i int, s *goquery.Selection) {
		p, err := m.parseProduct(s)
		if err != nil {
			return
		}

		p.Category = categoryName
		instruments[i] = p
	})

	return ProductResponse{Products: instruments}, nil
}

func (m MusikProduktiv) parseProduct(s *goquery.Selection) (products.Product, error) {
	manufacturer, model := m.parseProductName(s.Find("b").First().Text())
	price, err := m.parsePrice(s.Find("i").Text())
	if err != nil {
		return products.Product{}, err
	}

	return products.Product{
		Retailer:     "Musik Produktiv",
		Manufacturer: manufacturer,
		Model:        model,
		Price:        price,
		IsAvailable:  !s.Find(".ampel").HasClass("zzz"),
		ProductURL:   s.Find("a").First().AttrOr("href", ""),
		ThumbnailURL: s.Find("img").First().AttrOr("src", ""),
	}, nil
}

func (m MusikProduktiv) parseProductName(productName string) (manufacturer, model string) {
	parts := strings.Split(productName, " ")

	return parts[0], strings.TrimPrefix(productName, parts[0]+" ")
}

func (m MusikProduktiv) parsePrice(price string) (float64, error) {
	re := regexp.MustCompile("[^0-9]")
	p := re.ReplaceAllString(price, "")

	fPrice, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return 0, err
	}

	return fPrice, nil
}
