package main

import (
	"github.com/chrismeh/lefty/pkg/products"
	"math"
	"net/http"
	"strconv"
)

func (a application) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filter := products.Filter{Page: 1, ProductsPerPage: 50}
	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = uint(p)
		}
	}

	prds, err := a.productStore.FindAll(filter)
	if err != nil {
		a.jsonError(w, "Internal server jsonError", http.StatusInternalServerError)
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
		a.jsonError(w, "Internal server jsonError", http.StatusInternalServerError)
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
