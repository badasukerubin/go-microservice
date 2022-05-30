package server

import (
	"context"
	"log"

	protos "github.com/badasukerubin/go-microservices/protos/product"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l log.Logger) *Product {
	return &Product{l: &l}
}

func (p *Product) GetProduct(ctx context.Context, rr *protos.ProductRequest) (*protos.ProductResponse, error) {
	p.l.Print("Handle GetRate ", "Id: ", rr.GetID())

	return &protos.ProductResponse{Product: 1}, nil
}
