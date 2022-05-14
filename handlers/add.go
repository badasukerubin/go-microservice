package handlers

import (
	"fmt"
	"net/http"

	"github.com/badasukerubin/go-microservices/data"
)

// swagger:route POST /products products addProduct
// Adds a product
// responses:
//  200: noContent

// AddProduct adds a product
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Print(p)
	// p.l.Println("Handle Post Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	fmt.Printf("Prod: %#v", prod)

	data.AddProduct(&prod)
	// p.l.Printf("Prod: %#v", prod)
}
