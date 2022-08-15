package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Coffee -
type Coffee struct {
	log       hclog.Logger
	telemetry *telemetry.Telemetry
	con       data.Connection
}

// NewCoffee
func NewCoffee(t *telemetry.Telemetry, l hclog.Logger, con data.Connection) *Coffee {
	t.AddMeasure("coffee.create_coffee")

	return &Coffee{l, t, con}
}

func (c *Coffee) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.log.Info("Handle Coffee")

	vars := mux.Vars(r)

	var coffeeID *int

	if vars["id"] != "" {
		cId, err := strconv.Atoi(vars["id"])
		if err != nil {
			c.log.Error("CoffeeID provided could not be converted to an integer", "error", err)
			http.Error(rw, "Unable to list ingredients", http.StatusInternalServerError)
			return
		}
		coffeeID = &cId
	}

	cofs, err := c.con.GetCoffees(coffeeID)
	if err != nil {
		c.log.Error("Unable to get products from database", "error", err)
		http.Error(rw, "Unable to list products", http.StatusInternalServerError)
		return
	}

	d, err := cofs.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert products to JSON", "error", err)
		http.Error(rw, "Unable to list products", http.StatusInternalServerError)
		return
	}

	rw.Write(d)
}

// CreateCoffee creates a new coffee
func (c *Coffee) CreateCoffee(_ int, rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("coffee.create_coffee")
	defer done()

	c.log.Info("Handle Coffee | CreateCoffee")

	body := model.Coffee{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("Unable to decode JSON", "error", err)
		http.Error(rw, "Unable to parse request body", http.StatusInternalServerError)
		return
	}

	coffee, err := c.con.CreateCoffee(body)
	if err != nil {
		c.log.Error("Unable to create new coffee", "error", err)
		http.Error(rw, fmt.Sprintf("Unable to create new coffee: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	d, err := coffee.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert coffee to JSON", "error", err)
		http.Error(rw, "Unable to create new coffee", http.StatusInternalServerError)
		return
	}

	rw.Write(d)
}
