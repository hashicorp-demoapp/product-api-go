package data

import (
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/stretchr/testify/mock"
)

type MockConnection struct {
	mock.Mock
}

func (c *MockConnection) IsConnected() (bool, error) {
	return true, nil
}

func (c *MockConnection) GetProducts() (model.Coffees, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Coffees); ok {
		return m, args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *MockConnection) GetIngredientsForCoffee(coffeeid int) (model.Ingredients, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Ingredients); ok {
		return m, args.Error(1)
	}

	return nil, args.Error(1)
}
