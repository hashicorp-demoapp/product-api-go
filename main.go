package main

import (
	"net/http"
	"time"
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
	db, err := retryDBUntilReady()
	if err != nil {
		logger.Error("Timeout waiting for database connection")
		os.Exit(1)
	}

	r := mux.NewRouter()
	
	healthHandler := handlers.NewHealth(t, logger, db)
	r.Handle("/health", healthHandler).Methods("GET")

	coffeeHandler := handlers.NewCoffee(db, logger)
	r.Handle("/coffee", coffeeHandler).Methods("GET")

	http.ListenAndServe(":9090", r)
}

// retryDBUntilReady keeps retrying the database connection
// when running the application on a scheduler it is possible that the app will come up before
// the database, this can cause the app to go into a CrashLoopBackoff cycle
func retryDBUntilReady() (data.Connection, error) {
	st := time.Now()
	dt := 1*time.Second // this should be an exponential backoff
	mt := 60*time.Second // max time to wait of the DB connection

	for {
		db, err := data.New(conf.DBConnection)
		if err == nil {
			return db, nil
		}
	
		logger.Error("Unable to connect to database", "error", err)
		
		// check if max time has elapsed
		if time.Now().Sub(st) > mt {
			return nil, err
		}

		// retry
		time.Sleep(dt)
	}
}

func configUpdated() {
	logger.Info("Config file changed")
}
