package handlers

import (
	"net/http"

	"github.com/badasukerubin/go-microservices/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

// GetProducts returns the products
func (p *Products) GetProducts(w http.ResponseWriter, h *http.Request) {
	// Get the products
	lp := data.GetProducts()

	// Serialize the data to JSON
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}
