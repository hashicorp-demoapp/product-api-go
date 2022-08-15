package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Ingredients -
type Ingredients struct {
	log       hclog.Logger
	telemetry *telemetry.Telemetry
	con       data.Connection
}

// NewIngredients -
func NewIngredients(t *telemetry.Telemetry, l hclog.Logger, con data.Connection) *Ingredients {
	t.AddMeasure("ingredients.create_coffee_ingredient")

	return &Ingredients{l, t, con}
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
	done := c.telemetry.NewTiming("ingredients.create_coffee_ingredient")
	defer done()

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
