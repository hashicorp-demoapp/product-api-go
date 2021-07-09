package model

import (
	"encoding/json"
	"io"
)

// Token defines a JWT in the database
type Token struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"-"`
	DeletedAt string `db:"deleted_at" json:"-"`
}

// FromJSON serializes data from json
func (t *Token) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(t)
}

// ToJSON converts the collection to json
func (t *Token) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}
