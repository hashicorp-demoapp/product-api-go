package data

import (
	"github.com/hashicorp-demoapp/product-api-go/data/model"

	//"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection interface {
	IsConnected() (bool, error)
	GetProducts() (model.Coffees, error)
	GetIngredientsForCoffee(int) (model.Ingredients, error)
}

type PostgresSQL struct {
	db *sqlx.DB
}

// New creates a new connection to the database
func New(connection string) (Connection, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &PostgresSQL{db}, nil
}

// IsConnected checks the connection to the database and returns an error if not connected
func (c *PostgresSQL) IsConnected() (bool, error) {
	err := c.db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetProducts returns all products from the database
func (c *PostgresSQL) GetProducts() (model.Coffees, error) {
	cos := model.Coffees{}

	err := c.db.Select(&cos, "SELECT * FROM coffees")
	if err != nil {
		return nil, err
	}

	// fetch the ingredients for each coffee
	for n, cof := range cos {
		i := []model.CoffeeIngredients{}
		err := c.db.Select(&i, "SELECT ingredient_id FROM coffee_ingredients WHERE coffee_id=$1", cof.ID)
		if err != nil {
			return nil, err
		}

		cos[n].Ingredients = i
	}

	return cos, nil
}

// GetIngredientsForCoffee get the ingredients for the given coffeeid
func (c *PostgresSQL) GetIngredientsForCoffee(coffeeid int) (model.Ingredients, error) {
	is := []model.Ingredient{}

	err := c.db.Select(&is,
		`SELECT ingredients.id, ingredients.name, coffee_ingredients.quantity, coffee_ingredients.unit FROM ingredients 
		 LEFT JOIN coffee_ingredients ON ingredients.id=coffee_ingredients.ingredient_id 
		 WHERE coffee_ingredients.coffee_id=$1 AND coffee_ingredients.deleted_at IS NULL`,
		coffeeid,
	)
	if err != nil {
		return nil, err
	}

	return is, nil
}
