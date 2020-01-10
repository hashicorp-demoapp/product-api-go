package model

import (
	"encoding/json"
	"database/sql"
)

// Coffees is a list of Coffee
type Coffees []Coffee

// FromJSON serializes data from json
func (c*Coffees) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

// ToJSON converts the collection to json
func (c*Coffees) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Coffee defines a coffee in the database
type Coffee struct {
	ID 			int                 `db:"id" json:"id"`
	Name        string              `db:"name" json:"name"`
	Price       float64             `db:"price" json:"price"`
	CreatedAt   string              `db:"created_at" json:"created_at"`
	UpdatedAt   string              `db:"updated_at" json:"updated_at"`
	DeletedAt   sql.NullString      `db:"deleted_at" json:"deleted_at,omitempty"`
	Ingredients []CoffeeIngredients `json:"ingredients"`
}

type CoffeeIngredients struct {
	ID 			 int `db:"id" json:"-"`
	CoffeeID     int `db:"coffee_id" json:"-"`
	IngredientID int `db:"ingredient_id" json:"ingredient_id"`
}
