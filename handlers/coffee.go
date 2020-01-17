package handlers

import (
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp/go-hclog"
	"net/http"
)

type Coffee struct {
	con data.Connection
	log hclog.Logger
}

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
