package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/annbelievable/nicholas_jasckson_microservice_guide/product-api/data"
	"github.com/annbelievable/nicholas_jasckson_microservice_guide/product-api/handlers"
	gohandlers "github.com/gorrila/handlers"
	"github.com/gorrila/mux"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()

	l := log.New(os.Stout, "products-api", log.LstdFlags)
	v := data.NewValidation()

	//create the handlers
	ph := handlers.NewProducts(l, v)

	//create a new serve mux and register the handlers
	sm := mux.NewRouter()

	//handlers for API
	getR := sm.Methods(http.MethodGet).SubRouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).SubRouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddleWareValidateProduct)

	postR := sm.Methods(http.MethodPost).SubRouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddleWareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).SubRouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	//handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("docs", sh)
	getR.Handle("swagger.yaml", http.FileServer(http.Dir("./")))

	//CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	//create a new server
	s := http.Server{
		Addr:         *bindAddress,
		Handler:      ch(sm),
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting new server: %s\n", err)
			os.Exit(1)
		}
	}()

	//trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	//Block until a signal is received
	sig := <-c
	log.Println("Got signal:", sig)

	//gracefully shutdown the server, waiting max 30 seconds for current opertations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}