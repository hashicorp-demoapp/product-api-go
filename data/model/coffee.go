package model

import "github.com/jinzhu/gorm"

import "encoding/json"

// Coffees is a list of Coffee
type Coffees []Coffee

// FromJSON serializes data from JSON
func (c*Coffees) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

// ToJSON converts the collection to JSON
func (c*Coffees) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Coffee defines a coffee in the database
type Coffee struct {
	gorm.Model
	Name         string      `json:"name"`
	Price        float64     `json:"price"`
	Ingredients []Ingredient `gorm:"many2many:coffee_ingredients;"` 
}

// Ingredient defines an ingredient in the database
type Ingredient struct {
	gorm.Model
	Name         string `json:"name"`
	Quantity     string `json:"quantity"`
}