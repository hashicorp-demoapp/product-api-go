package handlers

import (
	"net/http"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp-demoapp/product-api-go/data"
)

type Ingredients struct {
	con data.Connection
	log hclog.Logger
}

func NewIngredients(con data.Connection, l hclog.Logger) *Ingredients {
	return &Ingredients{con, l}
}

func (c*Ingredients) ServeHTTP(rw http.ResponseWriter, r*http.Request) {
	c.log.Info("Handle Coffee")

	ingredients, err := c.con.GetIngredientsForCoffee(1)
	if err != nil {
		c.log.Error("Unable to get ingredients from database", "error", err)
		http.Error(rw, "Unable to list ingredients", http.StatusInternalServerError)
	}

	d, err := ingredients.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert products to JSON", "error", err)
		http.Error(rw, "Unable to list products", http.StatusInternalServerError)
	}

	rw.Write(d)
}