package main

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/handlers"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
	"net/http"
)

// Config format for application
var Config struct {
}

func main() {

	l := hclog.Default()
	t := telemetry.New()

	healthHandler := handlers.NewHealth(t, l)

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler.Handle).Methods("GET")

	http.ListenAndServe(":9090", r)
}
