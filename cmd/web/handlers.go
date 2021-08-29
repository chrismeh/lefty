package main

import (
	"net/http"
)

func (a application) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prds, err := a.productStore.FindAll()
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
