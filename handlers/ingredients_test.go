package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func setupIngredientsHandler() (*Ingredients, *httptest.ResponseRecorder) {
	c := &data.MockConnection{}
	c.On("GetIngredientsForCoffee").Return(model.Ingredients{
		model.Ingredient{ID: 1, Name: "Coffee"},
		model.Ingredient{ID: 2, Name: "Milk"},
		model.Ingredient{ID: 2, Name: "Sugar"},
	}, nil)
	c.On("CreateCoffeeIngredient").Return(model.CoffeeIngredients{
		ID:           1,
		CoffeeID:     2,
		IngredientID: 3,
	}, nil)

	l := hclog.Default()

	return &Ingredients{c, l}, httptest.NewRecorder()
}

func TestCoffeeIDReturnsIngredients(t *testing.T) {
	c, rw := setupIngredientsHandler()
	r := httptest.NewRequest("GET", "/coffees/{id:[0-9]+}/ingredients", nil)
	// set coffeeID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	c.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Ingredients{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}

// CreateCoffeeIngredient - Tests success criteria
func TestCreateCoffeeIngredient(t *testing.T) {
	c, rw := setupIngredientsHandler()

	userID := 1
	r := httptest.NewRequest("POST", "/coffees/{id:[0-9]+}/ingredients", nil)

	rb := strings.NewReader(`{"id":1,"coffee_id":2, "ingredient_id": 3}`)
	r.Body = ioutil.NopCloser(rb)

	c.CreateCoffeeIngredient(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Coffee{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.NoError(t, err)
}
