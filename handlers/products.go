package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/badasukerubin/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, h *http.Request) {
	if h.Method == http.MethodGet {
		p.getProducts(w, h)
		return
	}

	if h.Method == http.MethodPost {
		p.addProduct(w, h)
		return
	}

	if h.Method == http.MethodPut {
		r := regexp.MustCompile(`/([0-9]+)`)
		g := r.FindAllStringSubmatch(h.URL.Path, -1)

		if len(g) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, w, h)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(w http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, h *http.Request) {
	fmt.Print(p)
	// p.l.Println("Handle Post Product")

	prod := &data.Product{}
	err := prod.FromJSON(h.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal hson", http.StatusBadRequest)
	}
	fmt.Printf("Prod: %#v", prod)

	data.AddProduct(prod)
	// p.l.Printf("Prod: %#v", prod)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, h *http.Request) {
	// p.l.Println("Handle Post Product")

	prod := &data.Product{}
	err := prod.FromJSON(h.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal hson", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err != data.ErrorProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
	// p.l.Printf("Prod: %#v", prod)
}
