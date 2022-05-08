package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/badasukerubin/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

type KeyProduct struct{}

func NewProducts(l *log.Logger) *Products {
	return &Products{}
}

func (p *Products) GetProducts(w http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Print(p)
	// p.l.Println("Handle Post Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	fmt.Printf("Prod: %#v", prod)

	data.AddProduct(&prod)
	// p.l.Printf("Prod: %#v", prod)
}

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
