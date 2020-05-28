package model

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngredientsDeserializeFromJSON(t *testing.T) {
	c := Coffees{}

	err := c.FromJSON(bytes.NewReader([]byte(ingredientsData)))
	assert.NoError(t, err)

	assert.Len(t, c, 3)
	assert.Equal(t, 0, c[0].ID)
	assert.Equal(t, 1, c[1].ID)
	assert.Equal(t, 2, c[2].ID)
}

func TestIngredientsSerializesToJSON(t *testing.T) {
	c := Ingredients{
		Ingredient{ID: 1, Name: "test", Quantity: 10, Unit: "ml"},
	}

	d, err := c.ToJSON()
	assert.NoError(t, err)

	id := make([]map[string]interface{}, 0)
	err = json.Unmarshal(d, &id)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), id[0]["id"])
	assert.Equal(t, "test", id[0]["name"])
	assert.Equal(t, float64(10), id[0]["quantity"])
	assert.Equal(t, "ml", id[0]["unit"])
}

var ingredientsData = `
[
   {
      "id":0,
      "name":"Espresso",
      "quantity":40,
      "unit":"ml"
   },
   {
      "id":1,
      "name":"Semi Skimmed Milk",
      "quantity":300,
      "unit":"ml"
   },
   {
      "id":2,
      "name":"Pumpkin Spice",
      "quantity":5,
      "unit":"g"
   }
]
`
