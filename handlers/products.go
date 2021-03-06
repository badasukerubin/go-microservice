// Product API
//
// Documentation for Product API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/badasukerubin/go-microservices/data"
)

// A list of products to return in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in our DB
	// in: body
	Body []data.Product
}

// swagger:parameters updateProduct deleteProduct
type productIDParameterWrapper struct {
	// The ID of the product to update
	// in: path
	// required: true
	ID int `json:"id"`
}

// No content to return in the response
// swagger:response noContent
type productNoContentWrapper struct {
}

type Products struct {
	l *log.Logger
}

type KeyProduct struct{}

func NewProducts(l *log.Logger) *Products {
	return &Products{}
}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, *prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
