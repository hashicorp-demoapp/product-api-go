package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
)

// Coffee -
type Coffee struct {
	con data.Connection
	log hclog.Logger
}

// NewCoffee
func NewCoffee(con data.Connection, l hclog.Logger) *Coffee {
	return &Coffee{con, l}
}

func (c *Coffee) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.log.Info("Handle Coffee")

	prods, err := c.con.GetProducts()
	if err != nil {
		c.log.Error("Unable to get products from database", "error", err)
		http.Error(rw, "Unable to list products", http.StatusInternalServerError)
	}

	d, err := prods.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert products to JSON", "error", err)
		http.Error(rw, "Unable to list products", http.StatusInternalServerError)
	}

	rw.Write(d)
}

// CreateCoffee creates a new coffee
func (c *Coffee) CreateCoffee(_ int, rw http.ResponseWriter, r *http.Request) {
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
	}

	rw.Write(d)
}
