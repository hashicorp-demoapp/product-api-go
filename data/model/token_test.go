package model

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenDeserializeFromJSON(t *testing.T) {
	to := Token{}

	err := to.FromJSON(bytes.NewReader([]byte(tokenData)))
	assert.NoError(t, err)

	assert.Equal(t, 1, to.ID)
	assert.Equal(t, 5, to.UserID)
}

func TestTokenSerializesToJSON(t *testing.T) {
	to := Token{ID: 1, UserID: 5}

	d, err := to.ToJSON()
	assert.NoError(t, err)

	tod := make(map[string]interface{}, 0)
	err = json.Unmarshal(d, &tod)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), tod["id"])
	assert.Equal(t, float64(5), tod["user_id"])
}

var tokenData = `
{
	"id": 1,
	"user_id": 5
}
`
