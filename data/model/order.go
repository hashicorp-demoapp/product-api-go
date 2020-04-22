package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Orders is a list of Order
type Orders []Order

// FromJSON serializes data from json
func (o *Orders) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(o)
}

// ToJSON converts the collection to json
func (o *Orders) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

// Order defines an order in the database
type Order struct {
	ID        int            `db:"id" json:"id,omitempty"`
	UserID    int            `db:"user_id" json:"-"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
	Items     []OrderItems   `json:"items,omitempty"`
}

// FromJSON serializes data from json
func (o *Order) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(o)
}

// ToJSON converts the collection to json
func (o *Order) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

// OrderItems is an item/quantity in an order
type OrderItems struct {
	ID        int            `db:"id" json:"-"`
	OrderID   int            `db:"order_id" json:"-"`
	CoffeeID  int            `db:"coffee_id" json:"-"`
	Coffee    Coffee         `json:"coffee,omitempty"`
	Quantity  int            `db:"quantity" json:"quantity,omitempty"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
}
