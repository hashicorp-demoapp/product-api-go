package model

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrdersDeserializeFromJSON(t *testing.T) {
	o := Orders{}

	err := o.FromJSON(bytes.NewReader([]byte(ordersData)))
	assert.NoError(t, err)

	assert.Len(t, o, 2)
	assert.Equal(t, 0, o[0].ID)
	assert.Equal(t, 1, o[1].ID)

	fois := o[0].Items
	assert.Equal(t, 0, fois[0].Coffee.ID)
	assert.Equal(t, 1, fois[0].Quantity)
}

func TestOrdersSerializesToJSON(t *testing.T) {
	o := Orders{
		Order{
			ID: 1,
			Items: []OrderItems{
				OrderItems{
					Coffee:   Coffee{ID: 0},
					Quantity: 1,
				},
				OrderItems{
					Coffee:   Coffee{ID: 1},
					Quantity: 2,
				},
			},
		},
	}

	d, err := o.ToJSON()
	assert.NoError(t, err)

	od := make([]map[string]interface{}, 0)
	err = json.Unmarshal(d, &od)
	assert.NoError(t, err)

	order := od[0]

	assert.Equal(t, float64(1), order["id"])

	orderItems := order["items"].([]interface{})

	fois := orderItems[0].(map[string]interface{})
	assert.Equal(t, float64(0), fois["coffee"].(map[string]interface{})["id"])
	assert.Equal(t, float64(1), fois["quantity"])

	sois := orderItems[1].(map[string]interface{})
	assert.Equal(t, float64(1), sois["coffee"].(map[string]interface{})["id"])
	assert.Equal(t, float64(2), sois["quantity"])

}

var ordersData = `
[
   {
		"id":0,
		"items": [
			{
				"coffee": {
					"id":0,
					"name":"Latte"
				},
				"quantity":1
			},
			{
				"coffee": {
					"id":1,
					"name":"Americano"
				},
				"quantity":2
			}
		]
   },
   {
		"id":1,
		"items": [
			{
				"coffee": {
					"id":0,
					"name":"Latte"
				},
				"quantity":2
			},
			{
				"coffee": {
					"id":1,
					"name":"Americano"
				},
				"quantity":3
			}
		]
   }
]
`
