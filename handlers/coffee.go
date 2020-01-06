package handlers

import (
	"net/http"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp-demoapp/product-api-go/data"
)

type Coffee struct {
	con data.Connection
	log hclog.Logger
}

func NewCoffee(con data.Connection, l hclog.Logger) *Coffee {
	return &Coffee{con, l}
}

func (c*Coffee) ServeHTTP(rw http.ResponseWriter, r*http.Request) {
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