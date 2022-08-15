package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Order -
type Order struct {
	log       hclog.Logger
	telemetry *telemetry.Telemetry
	con       data.Connection
}

// NewOrder -
func NewOrder(t *telemetry.Telemetry, l hclog.Logger, con data.Connection) *Order {
	t.AddMeasure("order.get_user_orders")
	t.AddMeasure("order.create_order")
	t.AddMeasure("order.get_user_order")
	t.AddMeasure("order.update_order")
	t.AddMeasure("order.delete_order")

	return &Order{l, t, con}
}

func (c *Order) ServeHTTP(userID int, rw http.ResponseWriter, r *http.Request) {
	c.log.Info("Handle Order | unknown", "path", r.URL.Path)
	http.NotFound(rw, r)
}

// GetUserOrders gets all user orders for a specific user
func (c *Order) GetUserOrders(userID int, rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("order.get_user_orders")
	defer done()

	c.log.Info("Handle Orders | GetUserOrders")

	orders, err := c.con.GetOrders(userID, nil)
	if err != nil {
		c.log.Error("Unable to get order from database", "error", err)
		http.Error(rw, "Unable to list orders", http.StatusInternalServerError)
		return
	}

	d, err := orders.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert orders to JSON", "error", err)
		http.Error(rw, "Unable to list orders", http.StatusInternalServerError)
		return
	}

	rw.Write(d)
}

// CreateOrder creates a new order
func (c *Order) CreateOrder(userID int, rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("order.create_order")
	defer done()

	c.log.Info("Handle Orders | CreateOrder")

	body := []model.OrderItems{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("Unable to decode JSON", "error", err)
		http.Error(rw, "Unable to parse request body", http.StatusInternalServerError)
		return
	}

	order, err := c.con.CreateOrder(userID, body)
	if err != nil {
		c.log.Error("Unable to create new order", "error", err)
		http.Error(rw, "Unable to create new order", http.StatusInternalServerError)
		return
	}

	d, err := order.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert order to JSON", "error", err)
		http.Error(rw, "Unable to create new order", http.StatusInternalServerError)
	}

	rw.Write(d)
}

// GetUserOrder gets a specific user order
func (c *Order) GetUserOrder(userID int, rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("order.get_user_order")
	defer done()

	c.log.Info("Handle Orders | GetUserOrder")

	vars := mux.Vars(r)

	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Error("orderID provided could not be converted to an integer", "error", err)
		http.Error(rw, "Unable to list order", http.StatusInternalServerError)
		return
	}

	orders, err := c.con.GetOrders(userID, &orderID)
	if err != nil {
		c.log.Error("Unable to get order from database", "error", err)
		http.Error(rw, "Unable to list order", http.StatusInternalServerError)
		return
	}

	order := model.Order{}

	if len(orders) > 0 {
		order = orders[0]
	}

	d, err := order.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert orders to JSON", "error", err)
		http.Error(rw, "Unable to list order", http.StatusInternalServerError)
		return
	}

	rw.Write(d)
}

// UpdateOrder updates an order
func (c *Order) UpdateOrder(userID int, rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("order.update_order")
	defer done()

	c.log.Info("Handle Orders | UpdateOrder")

	// Get orderID
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Error("orderID provided could not be converted to an integer", "error", err)
		http.Error(rw, "Unable to update order", http.StatusInternalServerError)
		return
	}

	body := []model.OrderItems{}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		c.log.Error("Unable to decode JSON", "error", err)
		http.Error(rw, "Unable to parse request body", http.StatusInternalServerError)
		return
	}

	order, err := c.con.UpdateOrder(userID, orderID, body)
	if err != nil {
		c.log.Error("Unable to create new order", "error", err)
		http.Error(rw, "Unable to update order", http.StatusInternalServerError)
		return
	}

	d, err := order.ToJSON()
	if err != nil {
		c.log.Error("Unable to convert order to JSON", "error", err)
		http.Error(rw, "Unable to update order", http.StatusInternalServerError)
	}

	rw.Write(d)
}

// DeleteOrder deletes a user order
func (c *Order) DeleteOrder(userID int, rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("order.delete_order")
	defer done()

	c.log.Info("Handle Orders | DeleteOrder")

	vars := mux.Vars(r)

	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Error("orderID provided could not be converted to an integer", "error", err)
		http.Error(rw, "Unable to delete order", http.StatusInternalServerError)
		return
	}

	err = c.con.DeleteOrder(userID, orderID)
	if err != nil {
		c.log.Error("Unable to delete order from database", "error", err)
		http.Error(rw, "Unable to delete order", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s", "Deleted order")
}
