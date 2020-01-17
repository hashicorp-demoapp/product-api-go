package client

import (
	"fmt"
	"net/http"

	"github.com/hashicorp-demoapp/product-api-go/data/model"
)

// HTTP contains all client details
type HTTP struct {
	baseURL string
}

// NewHTTP creates a new HTTP client
func NewHTTP(baseURL string) *HTTP {
	return &HTTP{baseURL}
}

// GetCoffees retrieves a list of coffees
func (h *HTTP) GetCoffees() ([]model.Coffee, error) {
	resp, err := http.Get(fmt.Sprintf("%s/coffees", h.baseURL))
	if err != nil {
		return nil, err
	}

	coffees := model.Coffees{}
	coffees.FromJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	return coffees, nil
}

// GetCoffee retrieves a single coffee
func (h *HTTP) GetCoffee(coffeeID int) (*model.Coffee, error) {
	resp, err := http.Get(fmt.Sprintf("%s/coffees/%d", h.baseURL, coffeeID))
	if err != nil {
		return nil, err
	}

	coffee := model.Coffee{}
	err = coffee.FromJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	return &coffee, nil
}

// GetIngredientsForCoffee retrieves a list of ingredients that go into a particular coffee
func (h *HTTP) GetIngredientsForCoffee(coffeeID int) ([]model.Ingredient, error) {
	resp, err := http.Get(fmt.Sprintf("%s/coffees/%d/ingredients", h.baseURL, coffeeID))
	if err != nil {
		return nil, err
	}

	ingredients := model.Ingredients{}
	err = ingredients.FromJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}
