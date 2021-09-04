package retailer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chrismeh/lefty/pkg/products"
	"regexp"
	"strconv"
	"strings"
)

type MusikProduktiv struct {
	http          httpGetter
	manufacturers []string
}

func NewMusikProduktiv(http httpGetter) *MusikProduktiv {
	return &MusikProduktiv{http: http}
}

func (m *MusikProduktiv) LoadProducts(category string, options RequestOptions) (ProductResponse, error) {
	resp, err := m.http.Get(m.buildURL(category, options))
	if err != nil {
		return ProductResponse{}, fmt.Errorf("could not fetch products from musik-produktiv.de: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("could not create goquery document from reader: %w", err)
	}
	categoryName := doc.Find("div.list_title h1").Text()

	currentPage, lastPage, err := m.parsePagination(doc)
	if err != nil {
		return ProductResponse{}, err
	}

	if uint(lastPage) < options.Page {
		return ProductResponse{}, fmt.Errorf("page %d out of bounds, last page is %d", options.Page, lastPage)
	}

	manufacturerNodes := doc.Find(".mp-filtermenu ul").First().Find("li span")
	m.manufacturers = make([]string, len(manufacturerNodes.Nodes))
	manufacturerNodes.Each(func(i int, s *goquery.Selection) {
		m.manufacturers[i] = s.Text()
	})

	instrumentNodes := doc.Find("ul.artgrid li")
	instruments := make([]products.Product, len(instrumentNodes.Nodes))
	instrumentNodes.Each(func(i int, s *goquery.Selection) {
		p, err := m.parseProduct(s)
		if err != nil {
			return
		}

		p.Category = categoryName
		instruments[i] = p
	})

	return ProductResponse{
		Products:    instruments,
		CurrentPage: uint(currentPage),
		LastPage:    uint(lastPage),
	}, nil
}

func (m *MusikProduktiv) Categories() []string {
	return []string{
		"e-gitarre-linkshaender",
		"westerngitarre-linkshaender",
		"linkshaender-konzertgitarren",
		"e-bass-linkshaender",
	}
}

func (m *MusikProduktiv) buildURL(category string, options RequestOptions) string {
	var page uint = 1

	if options.Page > 0 {
		page = options.Page
	}

	return fmt.Sprintf("https://www.musik-produktiv.de/%s/?p=%d", category, page)
}

func (m *MusikProduktiv) parseProduct(s *goquery.Selection) (products.Product, error) {
	manufacturer, model := m.parseProductName(s.Find("b").First().Text())
	price, err := m.parsePrice(s.Find("i").Text())
	if err != nil {
		return products.Product{}, err
	}

	return products.Product{
		Retailer:          "Musik Produktiv",
		Manufacturer:      manufacturer,
		Model:             model,
		Price:             price,
		IsAvailable:       !s.Find(".ampel").HasClass("zzz"),
		AvailabilityScore: m.parseAvailabilityScore(s),
		ProductURL:        s.Find("a").First().AttrOr("href", ""),
		ThumbnailURL:      s.Find("img").First().AttrOr("src", ""),
	}, nil
}

func (m *MusikProduktiv) parseProductName(productName string) (manufacturer, model string) {
	for _, man := range m.manufacturers {
		if strings.HasPrefix(productName, man) {
			return man, strings.TrimPrefix(productName, man+" ")
		}
	}

	parts := strings.Split(productName, " ")
	return parts[0], strings.TrimPrefix(productName, parts[0]+" ")
}

func (m *MusikProduktiv) parsePrice(price string) (float64, error) {
	re := regexp.MustCompile("[^0-9]")
	p := re.ReplaceAllString(price, "")

	fPrice, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return 0, err
	}

	return fPrice, nil
}

func (m *MusikProduktiv) parseAvailabilityScore(s *goquery.Selection) int {
	s = s.Find(".ampel")
	if s.HasClass("ggg") {
		return products.AvailabilityAvailable
	}
	if s.HasClass("ggy") {
		return products.AvailabilityWithinDays
	}
	if s.HasClass("gyy") {
		return products.AvailabilityWithinWeeks
	}

	return products.AvailabilityUnknown
}

func (m *MusikProduktiv) parsePagination(s *goquery.Document) (currentPage, lastPage int, err error) {
	pagination := s.Find(".list_page div")
	if len(pagination.Nodes) == 1 {
		return 1, 1, nil
	}

	cp := pagination.Find("div").Text()
	lp := pagination.Children().Last().Text()

	currentPage, err = strconv.Atoi(cp)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse current page from pagination: %w", err)
	}

	lastPage, err = strconv.Atoi(lp)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse last page from pagination: %w", err)
	}

	return currentPage, lastPage, nil
}
