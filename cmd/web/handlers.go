package main

import (
	"github.com/chrismeh/lefty/pkg/products"
	"net/http"
	"strconv"
)

func (a application) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filter := products.Filter{ProductsPerPage: 50}
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

	err = a.json(w, &prds)
	if err != nil {
		a.jsonError(w, "Internal server jsonError", http.StatusInternalServerError)
		return
	}
}
