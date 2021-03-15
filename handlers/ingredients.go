package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
)

// Ingredients -
type Ingredients struct {
	con data.Connection
	log hclog.Logger
}

// NewIngredients -
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

// CreateCoffeeIngredient creates a new coffee ingredient
func (c *Ingredients) CreateCoffeeIngredient(_ int, rw http.ResponseWriter, r *http.Request) {
	c.log.Info("Handle Coffee | CreateCoffeeIngredient")

	body := struct {
		CoffeeID     int    `json:"coffee_id"`
		IngredientID int    `json:"ingredient_id"`
		Quantity     int    `json:"quantity"`
		Unit         string `json:"unit"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("Unable to decode JSON", "error", err)
		http.Error(rw, "Unable to parse request body", http.StatusInternalServerError)
		return
	}

	coffeeIngredient, err := c.con.UpsertCoffeeIngredient(
		model.Coffee{ID: body.CoffeeID},
		model.Ingredient{
			ID:       body.IngredientID,
			Quantity: body.Quantity,
			Unit:     body.Unit,
		})
	if err != nil {
		c.log.Error("Unable to create new coffeeIngredient", "error", err)
		http.Error(rw, "Unable to create new coffeeIngredient", http.StatusInternalServerError)
		return
	}

	d, err := coffeeIngredient.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert coffeeIngredient to JSON", "error", err)
		http.Error(rw, "Unable to create new coffeeIngredient", http.StatusInternalServerError)
	}

	rw.Write(d)
}
