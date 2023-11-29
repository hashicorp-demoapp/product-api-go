package main

import (
	"math"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp-demoapp/go-hckit"
	"github.com/nicholasjackson/env"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/config"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/handlers"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Config format for application
type Config struct {
	DBConnection           string  `json:"db_connection"`
	BindAddress            string  `json:"bind_address"`
	MetricsAddress         string  `json:"metrics_address"`
	MaxRetries             int     `json:"max_retries"`
	BackoffExponentialBase float64 `json:"backoff_exponential_base"`
	ErrorRate              int     `json:"error_rate"`
}

var conf *Config
var logger hclog.Logger

var configFile = env.String("CONFIG_FILE", false, "./conf.json", "Path to JSON encoded config file")
var dbConnection = env.String("DB_CONNECTION", false, "", "db connection string")
var bindAddress = env.String("BIND_ADDRESS", false, "", "Bind address")
var metricsAddress = env.String("METRICS_ADDRESS", false, "", "Metrics address")
var maxRetries = env.Int("MAX_RETRIES", false, 60, "Maximum number of connection retries")
var backoffExponentialBase = env.Float64("BACKOFF_EXPONENTIAL_BASE", false, 1, "Exponential base number to calculate the backoff")
var errorRate = env.Int("ERROR_RATE", false, 0, "Integer between 0 and 100 that represents how often the service should return 500 error code")

const jwtSecret = "test"

func main() {
	logger = hclog.Default()

	err := env.Parse()
	if err != nil {
		logger.Error("Error parsing flags", "error", err)
		os.Exit(1)
	}

	closer, err := hckit.InitGlobalTracer("product-api")
	if err != nil {
		logger.Error("Unable to initialize Tracer", "error", err)
		os.Exit(1)
	}
	defer closer.Close()

	conf = &Config{
		DBConnection:           *dbConnection,
		BindAddress:            *bindAddress,
		MetricsAddress:         *metricsAddress,
		MaxRetries:             *maxRetries,
		BackoffExponentialBase: *backoffExponentialBase,
		ErrorRate:              *errorRate,
	}

	// load the config, unless provided by env
	if conf.DBConnection == "" || conf.BindAddress == "" {
		c, err := config.New(*configFile, conf, configUpdated)
		if err != nil {
			logger.Error("Unable to load config file", "error", err)
			os.Exit(1)
		}
		defer c.Close()
	}

	// configure the telemetry
	t := telemetry.New(conf.MetricsAddress)

	// load the db connection
	db, err := retryDBUntilReady()
	if err != nil {
		logger.Error("Timeout waiting for database connection")
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.Use(hckit.TracingMiddleware)

	// Enable CORS for all hosts
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "content-type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}).Handler)

	authMiddleware := handlers.NewAuthMiddleware(db, logger)
	errorMiddleware := handlers.NewErrorMiddleware(conf.ErrorRate, logger)
	r.Use(errorMiddleware.Middleware)

	healthHandler := handlers.NewHealth(t, logger, db)
	r.Handle("/health", healthHandler).Methods("GET")
	r.HandleFunc("/health/livez", healthHandler.Liveness).Methods("GET")
	r.HandleFunc("/health/readyz", healthHandler.Readiness).Methods("GET")

	coffeeHandler := handlers.NewCoffee(db, logger)
	r.Handle("/coffees", coffeeHandler).Methods("GET")
	r.Handle("/coffees/{id:[0-9]+}", coffeeHandler).Methods("GET")
	r.Handle("/coffees", authMiddleware.IsAuthorized(coffeeHandler.CreateCoffee)).Methods("POST")

	ingredientsHandler := handlers.NewIngredients(db, logger)
	r.Handle("/coffees/{id:[0-9]+}/ingredients", ingredientsHandler).Methods("GET")
	r.Handle("/coffees/{id:[0-9]+}/ingredients", authMiddleware.IsAuthorized(ingredientsHandler.CreateCoffeeIngredient)).Methods("POST")

	userHandler := handlers.NewUser(db, logger)
	r.HandleFunc("/signup", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/signin", userHandler.SignIn).Methods("POST")
	r.HandleFunc("/signout", userHandler.SignOut).Methods("POST")

	orderHandler := handlers.NewOrder(db, logger)
	r.Handle("/orders", authMiddleware.IsAuthorized(orderHandler.GetUserOrders)).Methods("GET")
	r.Handle("/orders", authMiddleware.IsAuthorized(orderHandler.CreateOrder)).Methods("POST")
	r.Handle("/orders/{id:[0-9]+}", authMiddleware.IsAuthorized(orderHandler.GetUserOrder)).Methods("GET")
	r.Handle("/orders/{id:[0-9]+}", authMiddleware.IsAuthorized(orderHandler.UpdateOrder)).Methods("PUT")
	r.Handle("/orders/{id:[0-9]+}", authMiddleware.IsAuthorized(orderHandler.DeleteOrder)).Methods("DELETE")

	logger.Info("Starting service", "bind", conf.BindAddress, "metrics", conf.MetricsAddress)
	err = http.ListenAndServe(conf.BindAddress, r)
	if err != nil {
		logger.Error("Unable to start server", "bind", conf.BindAddress, "error", err)
	}
}

// retryDBUntilReady keeps retrying the database connection
// when running the application on a scheduler it is possible that the app will come up before
// the database, this can cause the app to go into a CrashLoopBackoff cycle
func retryDBUntilReady() (data.Connection, error) {
	maxRetries := conf.MaxRetries
	backoffExponentialBase := conf.BackoffExponentialBase
	dt := 0

	retries := 0
	backoff := time.Duration(0) // backoff before attempting to conection

	for {
		db, err := data.New(conf.DBConnection)
		if err == nil {
			return db, nil
		}

		logger.Error("Unable to connect to database", "error", err)

		// check if current retry reaches the max number of allowed retries
		if retries > maxRetries {
			return nil, err
		}

		// retry
		retries++
		dt = int(math.Pow(backoffExponentialBase, float64(retries)))
		backoff = time.Duration(dt) * time.Second
		time.Sleep(backoff)
	}
}

func configUpdated() {
	logger.Info("Config file changed")
}
