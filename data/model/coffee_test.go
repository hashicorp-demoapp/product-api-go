package model

import (
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"testing"
)

func TestCoffeesDeserializeFromJSON(t*testing.T) {
	c := Coffees{}

	err := c.FromJSON([]byte(coffeesData))
	assert.NoError(t, err)

	assert.Len(t, c, 2)
	assert.Equal(t, 1, c[0].ID)
	assert.Equal(t, 2, c[1].ID)
}

func TestCoffeesSerializesToJSON(t*testing.T) {
	c := Coffees{
		Coffee{ID: 1, Name: "test", Price: 120.12},
	}

	d, err := c.ToJSON()
	assert.NoError(t, err)

	cd := make([]map[string]interface{}, 0)
	err = json.Unmarshal(d, &cd)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), cd[0]["id"])
	assert.Equal(t, "test", cd[0]["name"])
	assert.Equal(t, float64(120.12), cd[0]["price"])
}

var coffeesData = `
[
	{
		"id": 1,
		"name": "Latte",
		"price": 50.0
	},
	{
		"id": 2,
		"name": "Americano",
		"price": 30.0
	}
]
`