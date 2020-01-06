package model

import "github.com/jinzhu/gorm"

import "encoding/json"

type Coffees []Coffee

// FromJSON serializes data from json
func (c*Coffees) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c*Coffees) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

type Coffee struct {
	gorm.Model
	ID		     int `json:"id"`
	Name         string `json:"name"`
	Price        float64 `json:"price"`
}