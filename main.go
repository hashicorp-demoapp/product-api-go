package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	DBConnection   string `json:"db_connection"`
	BindAddress    string `json:"bind_address"`
	MetricsAddress string `json:"metrics_address"`
}

var conf *Config
var logger hclog.Logger

var configFile = env.String("CONFIG_FILE", false, "./conf.json", "Path to JSON encoded config file")

const jwtSecret = "test"

func main() {
	logger = hclog.Default()

	err := env.Parse()
	if err != nil {
		logger.Error("Error parsing flags", "error", err)
		os.Exit(1)
	}

	conf = &Config{}

	// load the config
	c, err := config.New(*configFile, conf, configUpdated)
	if err != nil {
		logger.Error("Unable to load config file", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// configure the telemetry
	t := telemetry.New(conf.MetricsAddress)

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
	r.Handle("/coffees", coffeeHandler).Methods("GET")

	ingredientsHandler := handlers.NewIngredients(db, logger)
	r.Handle("/coffees/{id:[0-9]+}/ingredients", ingredientsHandler).Methods("GET")

	userHandler := handlers.NewUser(db, logger)
	r.HandleFunc("/signup", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/signin", userHandler.SignIn).Methods("POST")

	orderHandler := handlers.NewOrder(db, logger)
	r.Handle("/orders", isAuthorizedMiddleware(orderHandler.GetUserOrders)).Methods("GET")
	r.Handle("/orders", isAuthorizedMiddleware(orderHandler.CreateOrder)).Methods("POST")
	r.Handle("/orders/{id:[0-9]+}", isAuthorizedMiddleware(orderHandler.GetUserOrder)).Methods("GET")
	r.Handle("/orders/{id:[0-9]+}", isAuthorizedMiddleware(orderHandler.UpdateOrder)).Methods("PUT")
	r.Handle("/orders/{id:[0-9]+}", isAuthorizedMiddleware(orderHandler.DeleteOrder)).Methods("DELETE")

	logger.Info("Starting service", "bind", conf.BindAddress, "metrics", conf.MetricsAddress)
	err = http.ListenAndServe(conf.BindAddress, r)
	if err != nil {
		logger.Error("Unable to start server", "bind", conf.BindAddress, "error", err)
	}
}

// isAuthorizedMiddleware
func isAuthorizedMiddleware(next func(userID int, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("Unable to parse JWT token", "path", r.URL.Path)
				http.Error(w, "Unauthorized0", http.StatusUnauthorized)
				return nil, nil
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			logger.Error("Unauthorized", "error", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["user_id"].(float64))
			next(userID, w, r)
			return
		}
	})
}

// retryDBUntilReady keeps retrying the database connection
// when running the application on a scheduler it is possible that the app will come up before
// the database, this can cause the app to go into a CrashLoopBackoff cycle
func retryDBUntilReady() (data.Connection, error) {
	st := time.Now()
	dt := 1 * time.Second  // this should be an exponential backoff
	mt := 60 * time.Second // max time to wait of the DB connection

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
