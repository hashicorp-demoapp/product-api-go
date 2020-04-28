package handlers

import (
	"encoding/json"
	"errors"
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

func setupOrderHandler(t *testing.T) (*Order, *httptest.ResponseRecorder) {
	c := &data.MockConnection{}

	testOrder := model.Order{
		ID: 1,
		Items: []model.OrderItems{
			model.OrderItems{
				Coffee: model.Coffee{
					ID:   1,
					Name: "Latte",
				},
				Quantity: 1,
			},
			model.OrderItems{
				Coffee: model.Coffee{
					ID:   2,
					Name: "Mocha",
				},
				Quantity: 4,
			},
		},
	}

	c.On("GetOrders").Return(model.Orders{testOrder}, nil)
	c.On("CreateOrder").Return(testOrder, nil)
	c.On("UpdateOrder").Return(testOrder, nil)
	c.On("DeleteOrder").Return(nil)

	l := hclog.Default()

	return &Order{c, l}, httptest.NewRecorder()
}

func setupFailedOrderHandler(t *testing.T) (*Order, *httptest.ResponseRecorder) {
	c := &data.MockConnection{}

	c.On("GetOrders").Return(nil, errors.New("Unable to retrieve order"))
	c.On("CreateOrder").Return(nil, errors.New("Unable to create order"))
	c.On("UpdateOrder").Return(nil, errors.New("Unable to update order"))
	c.On("DeleteOrder").Return(errors.New("Unable to delete order"))

	l := hclog.Default()

	return &Order{c, l}, httptest.NewRecorder()
}

// TestReturnsOrders - Tests success criteria
func TestReturnsOrders(t *testing.T) {
	c, rw := setupOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("GET", "/orders", nil)

	c.GetUserOrders(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Orders{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.NoError(t, err)
}

// TestUnableToReturnOrders - Tests failure criteria
func TestUnableToReturnOrders(t *testing.T) {
	c, rw := setupFailedOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("GET", "/orders", nil)

	c.GetUserOrders(userID, rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)

	bd := model.Orders{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.Error(t, err)
	assert.Equal(t, "Unable to list orders\n", string(rw.Body.Bytes()))
}

// TestCreateOrder - Tests success criteria
func TestCreateOrder(t *testing.T) {
	c, rw := setupOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("POST", "/orders", nil)

	rb := strings.NewReader(`[{"coffee":{"id":1,"name":"Latte"},"quantity":2},{"coffee":{"id":2,"name":"Americano"},"quantity":3}]`)
	r.Body = ioutil.NopCloser(rb)

	c.CreateOrder(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Order{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.NoError(t, err)
}

// TestUnableToCreateOrder - Tests failure criteria
func TestUnableToCreateOrder(t *testing.T) {
	c, rw := setupFailedOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("POST", "/orders", nil)

	rb := strings.NewReader(`[{"coffee":{"id":1,"name":"Latte"},"quantity":2},{"coffee":{"id":2,"name":"Americano"},"quantity":3}]`)
	r.Body = ioutil.NopCloser(rb)

	c.CreateOrder(userID, rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)

	bd := model.Order{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.Error(t, err)
	assert.Equal(t, "Unable to create new order\n", string(rw.Body.Bytes()))
}

// TestReturnSpecificOrder - Tests success criteria
func TestReturnSpecificOrder(t *testing.T) {
	c, rw := setupOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("GET", "/orders/{id:[0-9]+}", nil)

	// set orderID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	c.GetUserOrder(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Order{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.NoError(t, err)
}

// TestUnableToReturnSpecificOrder - Tests failure criteria
func TestUnableToReturnSpecificOrder(t *testing.T) {
	c, rw := setupFailedOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("GET", "/orders/{id:[0-9]+}", nil)

	// set orderID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	c.GetUserOrder(userID, rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)

	bd := model.Order{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.Error(t, err)
	assert.Equal(t, "Unable to list order\n", string(rw.Body.Bytes()))
}

// TestUpdateOrder - Tests success criteria
func TestUpdateOrder(t *testing.T) {
	c, rw := setupOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("PUT", "/orders/{id:[0-9]+}", nil)

	// set orderID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	rb := strings.NewReader(`[{"coffee":{"id":1,"name":"Latte"},"quantity":2},{"coffee":{"id":2,"name":"Americano"},"quantity":3}]`)
	r.Body = ioutil.NopCloser(rb)

	c.UpdateOrder(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.Order{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.NoError(t, err)
}

// TestUnableToUpdateOrder - Tests failure criteria
func TestUnableToUpdateOrder(t *testing.T) {
	c, rw := setupFailedOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("PUT", "/orders/{id:[0-9]+}", nil)

	// set orderID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	rb := strings.NewReader(`[{"coffee":{"id":1,"name":"Latte"},"quantity":2},{"coffee":{"id":2,"name":"Americano"},"quantity":3}]`)
	r.Body = ioutil.NopCloser(rb)

	c.UpdateOrder(userID, rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)

	bd := model.Order{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)

	assert.Error(t, err)
	assert.Equal(t, "Unable to update order\n", string(rw.Body.Bytes()))
}

// TestDelete - Tests success criteria
func TestDelete(t *testing.T) {
	c, rw := setupOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("DELETE", "/orders/{id:[0-9]+}", nil)

	// set orderID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	c.DeleteOrder(userID, rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "Deleted order", string(rw.Body.Bytes()))
}

// TestUnableToDelete - Tests failure criteria
func TestUnableToDelete(t *testing.T) {
	c, rw := setupFailedOrderHandler(t)

	userID := 1
	r := httptest.NewRequest("DELETE", "/orders/{id:[0-9]+}", nil)

	// set orderID to 1
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	c.DeleteOrder(userID, rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	assert.Equal(t, "Unable to delete order\n", string(rw.Body.Bytes()))
}
