package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Ingredients is a collection of Ingredient
type Ingredients []Ingredient

// FromJSON serializes data from json
func (c *Ingredients) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(c)
}

// ToJSON converts the collection to json
func (c *Ingredients) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Ingredient defines an ingredient in the database
type Ingredient struct {
	ID        int            `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Quantity  int            `db:"quantity" json:"quantity"`
	Unit      string         `db:"unit" json:"unit"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
}
