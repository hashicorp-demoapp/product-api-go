package handlers

import (
	"net/http"
	"strconv"

	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp/go-hclog"

	"github.com/gorilla/mux"
)

type Ingredients struct {
	con data.Connection
	log hclog.Logger
}

func NewIngredients(con data.Connection, l hclog.Logger) *Ingredients {
	return &Ingredients{con, l}
}

func (c *Ingredients) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.log.Info("Handle Coffee Ingredients")

	vars := mux.Vars(r)

	coffeeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Error("CoffeeID provided could not be converted to an integer", "error", err)
		http.Error(rw, "Unable to list ingredients", http.StatusInternalServerError)
	}

	ingredients, err := c.con.GetIngredientsForCoffee(coffeeID)
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
