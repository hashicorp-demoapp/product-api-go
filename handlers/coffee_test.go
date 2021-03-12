package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func setupCoffeeHandler() (*Coffee, *httptest.ResponseRecorder) {
	c := &data.MockConnection{}
	c.On("GetProducts").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)
	c.On("CreateCoffee").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)

	l := hclog.Default()

	return &Coffee{c, l}, httptest.NewRecorder()
}

func TestCoffeeReturnsProducts(t *testing.T) {
	c, rw := setupCoffeeHandler()
	r := httptest.NewRequest("GET", "/coffees", nil)
	c.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Coffees{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}

// TestCreateCoffee - Tests success criteria
func TestCreateCoffee(t *testing.T) {
	c, rw := setupCoffeeHandler()

	userID := 1
	r := httptest.NewRequest("POST", "/coffees", nil)

	rb := strings.NewReader(`{"coffee":{"id":1,"name":"Latte"}}`)
	r.Body = ioutil.NopCloser(rb)

	c.CreateCoffee(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Coffee{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.NoError(t, err)
}
