package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func setupIngredientsHandler(t *testing.T) (*Ingredients, *httptest.ResponseRecorder, *http.Request) {
	c := &data.MockConnection{}
	c.On("GetIngredientsForCoffee").Return(model.Ingredients{
		model.Ingredient{ID: 1, Name: "Coffee"},
		model.Ingredient{ID: 2, Name: "Milk"},
		model.Ingredient{ID: 2, Name: "Sugar"},
	}, nil)

	l := hclog.Default()

	return &Ingredients{c, l}, httptest.NewRecorder(), httptest.NewRequest("GET", "/coffees/{id:[0-9]+}/ingredients", nil)
}

func TestCoffeeIDReturnsIngredients(t *testing.T) {
	c, rw, r := setupIngredientsHandler(t)

	// set coffeeID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	c.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Ingredients{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}
