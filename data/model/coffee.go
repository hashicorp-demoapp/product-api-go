package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Coffees is a list of Coffee
type Coffees []Coffee

// FromJSON serializes data from json
func (c *Coffees) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(c)
}

// ToJSON converts the collection to json
func (c *Coffees) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Coffee defines a coffee in the database
type Coffee struct {
	ID          int                 `db:"id" json:"id"`
	Name        string              `db:"name" json:"name"`
	Teaser      string              `db:"teaser" json:"teaser"`
	Description string              `db:"description" json:"description"`
	Price       float64             `db:"price" json:"price"`
	Image       string              `db:"image" json:"image"`
	CreatedAt   string              `db:"created_at" json:"-"`
	UpdatedAt   string              `db:"updated_at" json:"-"`
	DeletedAt   sql.NullString      `db:"deleted_at" json:"-"`
	Ingredients []CoffeeIngredients `json:"ingredients"`
}

func (c *Coffee) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(c)
}

// ToJSON converts the collection to json
func (c *Coffee) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

type CoffeeIngredients struct {
	ID           int            `db:"id" json:"-"`
	CoffeeID     int            `db:"coffee_id" json:"-"`
	IngredientID int            `db:"ingredient_id" json:"ingredient_id"`
	CreatedAt    string         `db:"created_at" json:"-"`
	UpdatedAt    string         `db:"updated_at" json:"-"`
	DeletedAt    sql.NullString `db:"deleted_at" json:"-"`
}
