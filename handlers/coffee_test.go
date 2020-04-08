package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func setupCoffeeHandler(t *testing.T) (*Coffee, *httptest.ResponseRecorder, *http.Request) {
	c := &data.MockConnection{}
	c.On("GetProducts").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)

	l := hclog.Default()

	return &Coffee{c, l}, httptest.NewRecorder(), httptest.NewRequest("GET", "/coffees", nil)
}

func TestCoffeeReturnsProducts(t *testing.T) {
	c, rw, r := setupCoffeeHandler(t)

	c.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Coffees{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}
