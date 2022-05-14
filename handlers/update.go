package handlers

import (
	"net/http"
	"strconv"

	"github.com/badasukerubin/go-microservices/data"
	"github.com/gorilla/mux"
)

// swagger:route PUT /products/{id} products updateProduct
// Returns a list of products
// responses:
//  200: noContent

// UpdateProduct updates a product
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// p.l.Println("Handle Post Product")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
	// p.l.Printf("Prod: %#v", prod)
}
