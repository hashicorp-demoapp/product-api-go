package main

import (
	"net/http"
	"os"

	"github.com/nicholasjackson/env"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/config"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/handlers"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Config format for application
type Config struct {
	DBConnection string `json:"db_connection"`
}

var conf *Config
var logger hclog.Logger

var configFile = env.String("CONFIG_FILE", true, "", "Path to JSON encoded config file")

func main() {
	logger = hclog.Default()

	err := env.Parse()
	if err != nil {
		logger.Error("Error parsing flags", "error", err)
		os.Exit(1)
	}

	conf = &Config{}

	t := telemetry.New()

	// load the config
	c, err := config.New(*configFile, conf, configUpdated)
	if err != nil {
		logger.Error("Unable to load config file", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// load the db connection
	db, err := data.New(conf.DBConnection)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	
	healthHandler := handlers.NewHealth(t, logger, db)
	r.Handle("/health", healthHandler).Methods("GET")

	coffeeHandler := handlers.NewCoffee(db, logger)
	r.Handle("/coffee", coffeeHandler).Methods("GET")

	http.ListenAndServe(":9090", r)
}

func configUpdated() {
	logger.Info("Config file changed")
}
