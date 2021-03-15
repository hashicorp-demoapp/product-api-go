package data

import (
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/stretchr/testify/mock"
)

type MockConnection struct {
	mock.Mock
}

// IsConnected -
func (c *MockConnection) IsConnected() (bool, error) {
	return true, nil
}

// GetProducts -
func (c *MockConnection) GetProducts() (model.Coffees, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Coffees); ok {
		return m, args.Error(1)
	}

	return nil, args.Error(1)
}

// GetIngredientsForCoffee -
func (c *MockConnection) GetIngredientsForCoffee(coffeeid int) (model.Ingredients, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Ingredients); ok {
		return m, args.Error(1)
	}

	return nil, args.Error(1)
}

// CreateUser -
func (c *MockConnection) CreateUser(username string, password string) (model.User, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.User); ok {
		return m, args.Error(1)
	}

	return model.User{}, args.Error(1)
}

// AuthUser -
func (c *MockConnection) AuthUser(username string, password string) (model.User, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.User); ok {
		return m, args.Error(1)
	}

	return model.User{}, args.Error(1)
}

// GetOrders -
func (c *MockConnection) GetOrders(userID int, orderID *int) (model.Orders, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Orders); ok {
		return m, args.Error(1)
	}

	return nil, args.Error(1)
}

// CreateOrder -
func (c *MockConnection) CreateOrder(userID int, orderItems []model.OrderItems) (model.Order, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Order); ok {
		return m, args.Error(1)
	}

	return model.Order{}, args.Error(1)
}

// UpdateOrder -
func (c *MockConnection) UpdateOrder(userID int, orderID int, orderItems []model.OrderItems) (model.Order, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Order); ok {
		return m, args.Error(1)
	}

	return model.Order{}, args.Error(1)
}

// DeleteOrder -
func (c *MockConnection) DeleteOrder(userID int, orderID int) error {
	args := c.Called()

	if err, ok := args.Get(0).(error); ok {
		return err
	}

	return nil
}

// CreateCoffee creates a new coffee type
func (c *MockConnection) CreateCoffee(coffee model.Coffee) (model.Coffee, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.Coffee); ok {
		return m, args.Error(1)
	}

	return model.Coffee{}, args.Error(1)
}

// CreateCoffeeIngredient creates a new coffee ingredient type
func (c *MockConnection) UpsertCoffeeIngredient(coffee model.Coffee, ingredient model.Ingredient) (model.CoffeeIngredient, error) {
	args := c.Called()

	if m, ok := args.Get(0).(model.CoffeeIngredient); ok {
		return m, args.Error(1)
	}

	return model.CoffeeIngredient{}, args.Error(1)
}
