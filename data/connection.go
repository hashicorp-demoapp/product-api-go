package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Connection interface {
	IsConnected() (bool, error)
}

type PostgresSQL struct {
	db *sql.DB
}

func New(connection string) (Connection, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &PostgresSQL{db}, nil
}

// IsConnection check the connection to the database and returns an error if not connected
func (c *PostgresSQL) IsConnected() (bool, error) {
	err := c.db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}
