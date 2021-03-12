package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/handlers"
	"github.com/hashicorp/go-hclog"
)

func (api *apiFeature) initHandlers() {
	// Coffee
	mc := &data.MockConnection{}
	mc.On("GetProducts").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)
	mc.On("CreateCoffee").Return(model.Coffee{ID: 1, Name: "Test"}, nil)
	mc.On("CreateCoffeeIngredient").Return(model.CoffeeIngredients{ID: 1, CoffeeID: 1, IngredientID: 3}, nil)
	mc.On("GetIngredientsForCoffee").Return(model.Ingredients{
		model.Ingredient{ID: 1, Name: "Coffee"},
		model.Ingredient{ID: 2, Name: "Milk"},
		model.Ingredient{ID: 2, Name: "Sugar"},
	})
	// User
	mc.On("CreateUser").Return(model.User{ID: 1, Username: "User1"}, nil)
	mc.On("AuthUser").Return(model.User{ID: 1, Username: "User1"}, nil)
	// Order
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

	mc.On("GetOrders").Return(model.Orders{testOrder}, nil)
	mc.On("CreateOrder").Return(testOrder, nil)
	mc.On("UpdateOrder").Return(testOrder, nil)
	mc.On("DeleteOrder").Return(nil)

	l := hclog.Default()

	api.mc = mc
	api.hc = handlers.NewCoffee(mc, l)
	api.hu = handlers.NewUser(mc, l)
	api.ho = handlers.NewOrder(mc, l)
	api.hi = handlers.NewIngredients(mc, l)
}

func (api *apiFeature) initRouter(method, endpoint string, userID *string) error {
	if userID != nil {
		i, err := strconv.Atoi(*userID)
		if err != nil {
			return err
		}

		switch method {
		case http.MethodGet:
			if endpoint == "/orders" {
				api.ho.GetUserOrders(i, api.rw, api.r)
			} else {
				api.ho.GetUserOrder(i, api.rw, api.r)
			}
		case http.MethodPost:
			if endpoint == "/coffees/{id:[0-9]+}/ingredients" {
				api.hi.CreateCoffeeIngredient(i, api.rw, api.r)
				return nil
			}
			if endpoint == "/coffees" {
				api.hc.CreateCoffee(i, api.rw, api.r)
				return nil
			}
			api.ho.CreateOrder(i, api.rw, api.r)
		case http.MethodPut:
			api.ho.UpdateOrder(i, api.rw, api.r)
		case http.MethodDelete:
			api.ho.DeleteOrder(i, api.rw, api.r)
		}
		return nil
	}

	if strings.Contains(endpoint, "/coffees") {
		api.hc.ServeHTTP(api.rw, api.r)
		return nil
	}
	if strings.Contains(endpoint, "/signup") {
		api.hu.SignUp(api.rw, api.r)
		return nil
	}
	if strings.Contains(endpoint, "/signin") {
		api.hu.SignIn(api.rw, api.r)
	}
	return nil
}

func (api *apiFeature) theServerIsRunning() error {
	connected, err := api.mc.IsConnected()
	if err != nil {
		return err
	}
	if connected == false {
		return fmt.Errorf("Mock connection is not connected")
	}
	return nil
}

func (api *apiFeature) iMakeARequestTo(method, endpoint string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWhereIs(method, endpoint string, attribute, value string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	vars := map[string]string{attribute: value}
	api.r = mux.SetURLVars(api.r, vars)

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWithTheFollowingRequestBody(method, endpoint string, body *messages.PickleStepArgument_PickleDocString) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	rb := strings.NewReader(body.Content)
	api.r.Body = ioutil.NopCloser(rb)

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWhereMyUserIDIs(method, endpoint, userID string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	err := api.initRouter(method, endpoint, &userID)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWhereMyUserIDIsWithTheFollowingRequestBody(method, endpoint, userID string, body *messages.PickleStepArgument_PickleDocString) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	rb := strings.NewReader(body.Content)
	api.r.Body = ioutil.NopCloser(rb)

	err := api.initRouter(method, endpoint, &userID)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWhereIsAndMyUserIDIs(method, endpoint, attribute, value, userID string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	vars := map[string]string{attribute: value}
	api.r = mux.SetURLVars(api.r, vars)

	err := api.initRouter(method, endpoint, &userID)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWhereIsAndMyUserIDIsWithTheFollowingRequestBody(method, endpoint, attribute, value, userID string, body *messages.PickleStepArgument_PickleDocString) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	rb := strings.NewReader(body.Content)
	api.r.Body = ioutil.NopCloser(rb)

	vars := map[string]string{attribute: value}
	api.r = mux.SetURLVars(api.r, vars)

	err := api.initRouter(method, endpoint, &userID)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) aListOfProductsShouldBeReturned() error {
	bd := model.Coffees{}

	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) thatProductsIngredientsShouldBeReturned() error {
	bd := model.Ingredients{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) anOrderShouldBeReturned() error {
	bd := model.Order{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) aCoffeeShouldBeReturned() error {
	bd := model.Coffee{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(api.rw.Body.Bytes()))
	}
	return nil
}

func (api *apiFeature) aCoffeeIngredientShouldBeReturned() error {
	bd := model.CoffeeIngredients{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(api.rw.Body.Bytes()))
	}
	return nil
}

func (api *apiFeature) aListOfOrdersShouldBeReturned() error {
	bd := model.Orders{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) theAuthResponseShouldBeReturned() error {
	bd := handlers.AuthResponse{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) theResponseStatusShouldBe(statusCode string) error {
	switch statusCode {
	case "OK":
		if api.rw.Code != http.StatusOK {
			return fmt.Errorf("expected status code does not match actual, %v vs. %v", http.StatusOK, api.rw.Code)
		}
	default:
		return fmt.Errorf("Status Code is not valid, %s", statusCode)
	}
	return nil
}
