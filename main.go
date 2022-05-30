package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/badasukerubin/go-microservices/files"
	"github.com/badasukerubin/go-microservices/handlers"
	protos "github.com/badasukerubin/go-microservices/protos/product"
	server "github.com/badasukerubin/go-microservices/server"
)

func main() {
	l := log.New(os.Stdout, "microservice", log.LstdFlags)
	stor, err := files.NewLocal("./imagestore", 1024*1000*5)

	if err != nil {
		l.Fatal("Unable to create storage", err)
		os.Exit(1)
	}

	ph := handlers.NewProducts(l)
	fh := handlers.NewFiles(stor, l)
	mw := handlers.GzipHandler{}

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/products", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	ops := middleware.RedocOpts{SpecURL: "swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	pfh := sm.Methods(http.MethodPost).Subrouter()
	pfh.HandleFunc("/products/file", fh.UploadMultipart)

	pfh.Use(mw.GZipMiddleware)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	gs := grpc.NewServer()
	ps := server.NewProduct(*l)

	protos.RegisterProductServer(gs, ps)

	// Disable in prod
	reflection.Register(gs)

	ln, err := net.Listen("tcp", ":9092")
	if err != nil {
		l.Fatal("Unable to listen", err)
	}

	gs.Serve(ln)

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminat5e, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}
