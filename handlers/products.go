package handlers

import (
	"log"
	"net/http"

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

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(w http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}
