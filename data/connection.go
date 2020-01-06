package data

import (
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Connection interface {
	IsConnected() (bool, error)
	GetProducts() (model.Coffees, error)
}

type PostgresSQL struct {
	db *gorm.DB
}

// New creates a new connection to the database
func New(connection string) (Connection, error) {
	db, err := gorm.Open("postgres", connection)
	//db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &PostgresSQL{db}, nil
}

// IsConnected checks the connection to the database and returns an error if not connected
func (c *PostgresSQL) IsConnected() (bool, error) {
	err := c.db.DB().Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetProducts returns all products from the database
func (c*PostgresSQL) GetProducts() (model.Coffees, error) {
	cos := model.Coffees{}

	db := c.db.Find(&cos)
	if db.Error != nil {
		return nil, db.Error
	}

	return cos,nil
}
