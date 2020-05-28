package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/handlers"
	"github.com/hashicorp/go-hclog"
)

var runTest *bool = flag.Bool("run.test", false, "Should we run the tests")

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !*runTest {
		return
	}

	format := "progress"
	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: format,
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

type apiFeature struct {
	mc *data.MockConnection
	c  *handlers.Coffee
	rw *httptest.ResponseRecorder
	r  *http.Request
}

func (api *apiFeature) newCoffeeHandler() {
	mc := &data.MockConnection{}
	mc.On("GetProducts").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)
	mc.On("GetIngredientsForCoffee").Return(model.Ingredients{
		model.Ingredient{ID: 1, Name: "Coffee"},
		model.Ingredient{ID: 2, Name: "Milk"},
		model.Ingredient{ID: 2, Name: "Sugar"},
	})
	l := hclog.Default()

	api.mc = mc
	api.c = handlers.NewCoffee(mc, l)
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

	api.c.ServeHTTP(api.rw, api.r)

	return nil
}

func (api *apiFeature) iMakeARequestToWhereIs(method, endpoint string, attribute, value string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	vars := map[string]string{attribute: value}
	api.r = mux.SetURLVars(api.r, vars)

	api.c.ServeHTTP(api.rw, api.r)

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

func (api *apiFeature) thatProductsIngredientsShouldBeReturned() error {
	bd := model.Ingredients{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}

	api.newCoffeeHandler()

	s.Step(`^the server is running$`, api.theServerIsRunning)

	s.Step(`^I make a "([^"]*)" request to "([^"]*)"$`, api.iMakeARequestTo)
	s.Step(`^I make a "([^"]*)" request to "([^"]*)" where "([^"]*)" is "([^"]*)"$`, api.iMakeARequestToWhereIs)

	s.Step(`^a list of products should be returned$`, api.aListOfProductsShouldBeReturned)
	s.Step(`^the response status should be "([^"]*)"$`, api.theResponseStatusShouldBe)
	s.Step(`^a list of the product\'s ingredients should be returned$`, api.thatProductsIngredientsShouldBeReturned)
}
