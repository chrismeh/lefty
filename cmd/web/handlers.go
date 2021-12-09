package main

import (
	"github.com/chrismeh/lefty/pkg/products"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

func (a application) handleShowIndex(w http.ResponseWriter, _ *http.Request) {
	file, err := os.Open("./templates/index.html")
	if err != nil {
		a.errorLog.Printf("could not open template index.html: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		a.errorLog.Printf("could not read template index.html: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(content)
}

func (a application) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filter := products.Filter{
		Search:          r.URL.Query().Get("search"),
		OrderBy:         r.URL.Query().Get("order"),
		Page:            1,
		ProductsPerPage: 50,
	}
	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = uint(p)
		}
	}

	prds, err := a.productStore.FindAll(filter)
	if err != nil {
		a.jsonError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	count := a.productStore.Count(filter)
	lastPage := math.Ceil(float64(count) / float64(filter.ProductsPerPage))

	resp := response{
		Data: prds,
		Meta: meta{
			CurrentPage:  filter.Page,
			LastPage:     uint(lastPage),
			OverallCount: uint(count),
			Count:        uint(len(prds)),
		},
	}
	err = a.json(w, resp)
	if err != nil {
		a.jsonError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

type response struct {
	Data []products.Product `json:"data"`
	Meta meta               `json:"meta"`
}

type meta struct {
	CurrentPage  uint `json:"current_page"`
	LastPage     uint `json:"last_page"`
	OverallCount uint `json:"overall_count"`
	Count        uint `json:"count"`
}
