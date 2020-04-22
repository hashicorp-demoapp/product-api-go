package model

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserDeserializeFromJSON(t *testing.T) {
	u := User{}

	err := u.FromJSON(bytes.NewReader([]byte(userData)))
	assert.NoError(t, err)

	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "testUser", u.Username)
}

func TestUserSerializesToJSON(t *testing.T) {
	u := User{ID: 1, Username: "test"}

	d, err := u.ToJSON()
	assert.NoError(t, err)

	ud := make(map[string]interface{}, 0)
	err = json.Unmarshal(d, &ud)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), ud["id"])
	assert.Equal(t, "test", ud["username"])
}

var userData = `
{
	"id": 1,
	"username": "testUser"
}
`
