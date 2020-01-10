package model

import "encoding/json"

// Ingredients is a collection of Ingredient
type Ingredients [] Ingredient

// FromJSON serializes data from json
func (c*Ingredients) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

// ToJSON converts the collection to json
func (c*Ingredients) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Ingredient defines an ingredient in the database
type Ingredient struct {
	Name         string `json:"name"`
	Quantity     string `json:"quantity"`
}