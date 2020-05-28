package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// User defines a user in the database
type User struct {
	ID        int            `db:"id" json:"id"`
	Username  string         `db:"username" json:"username"`
	Password  string         `db:"password" json:"-"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
	Orders    []Orders       `json:"orders"`
}

// FromJSON serializes data from json
func (u *User) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(u)
}

// ToJSON converts the collection to json
func (u *User) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}
