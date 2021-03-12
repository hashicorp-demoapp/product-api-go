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
	CreateUser(string, string) (model.User, error)
	AuthUser(string, string) (model.User, error)
	GetOrders(int, *int) (model.Orders, error)
	CreateOrder(int, []model.OrderItems) (model.Order, error)
	UpdateOrder(int, int, []model.OrderItems) (model.Order, error)
	DeleteOrder(int, int) error
	CreateCoffee(model.Coffee) (model.Coffee, error)
	CreateCoffeeIngredient(model.Coffee, model.Ingredient) (model.CoffeeIngredients, error)
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

// CreateUser creates a new user
func (c *PostgresSQL) CreateUser(username string, password string) (model.User, error) {
	u := model.User{}

	rows, err := c.db.NamedQuery(
		`INSERT INTO users (username, password, created_at, updated_at) 
		VALUES(:username, crypt(:password, gen_salt('bf')), now(), now()) 
		RETURNING id, username;`, map[string]interface{}{
			"username": username,
			"password": password,
		})
	if err != nil {
		return u, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&u)
		if err != nil {
			return u, err
		}
	}

	return u, nil
}

// AuthUser checks whether username and password matches
func (c *PostgresSQL) AuthUser(username string, password string) (model.User, error) {
	us := []model.User{}

	err := c.db.Select(&us,
		`SELECT id, username FROM users 
		WHERE username = $1 AND password = crypt($2, password);`,
		username, password,
	)
	if err != nil {
		return us[0], err
	}

	return us[0], nil
}

// GetOrders returns orders from the database
func (c *PostgresSQL) GetOrders(userID int, orderID *int) (model.Orders, error) {
	orders := model.Orders{}

	if orderID != nil {
		err := c.db.Select(&orders,
			`SELECT * FROM orders WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL`,
			userID, orderID)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.db.Select(&orders,
			`SELECT * FROM orders WHERE user_id = $1 AND deleted_at IS NULL`,
			userID)
		if err != nil {
			return nil, err
		}
	}

	// fetch the coffee for each order
	for n, order := range orders {
		items := []model.OrderItems{}
		err := c.db.Select(&items,
			`SELECT * FROM order_items WHERE order_id=$1 AND deleted_at IS NULL`, order.ID)
		if err != nil {
			return nil, err
		}
		orders[n].Items = items

		for i, item := range items {
			coffee := model.Coffees{}
			err := c.db.Select(&coffee,
				`SELECT * FROM coffees WHERE id=$1 AND deleted_at IS NULL`, item.CoffeeID)
			if err != nil {
				return nil, err
			}

			if len(coffee) > 0 {
				orders[n].Items[i].Coffee = coffee[0]
			}
		}
	}

	return orders, nil
}

// CreateOrder creates a new order in the database
func (c *PostgresSQL) CreateOrder(userID int, orderItems []model.OrderItems) (model.Order, error) {
	tx := c.db.MustBegin()

	o := model.Order{}
	rows, err := tx.NamedQuery(
		`INSERT INTO orders (user_id, created_at, updated_at) 
		VALUES (:user_id, now(), now()) RETURNING id`, map[string]interface{}{
			"user_id": userID,
		})
	if err != nil {
		return o, err
	}
	if rows.Next() {
		err := rows.StructScan(&o)
		if err != nil {
			tx.Rollback()
			return o, err
		}
	}

	rows.Close()

	for _, item := range orderItems {
		_, err = tx.NamedExec(
			`INSERT INTO order_items (order_id, coffee_id, quantity, created_at, updated_at) 
			VALUES (:order_id, :coffee_id, :quantity, now(), now())`, map[string]interface{}{
				"order_id":  o.ID,
				"coffee_id": item.Coffee.ID,
				"quantity":  item.Quantity,
			})
		if err != nil {
			tx.Rollback()
			return o, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return o, err
	}

	orders, err := c.GetOrders(userID, &o.ID)
	if err != nil {
		return o, err
	}

	if len(orders) == 0 {
		return o, err
	}

	return orders[0], nil
}

// UpdateOrder updates an existing order in the database
func (c *PostgresSQL) UpdateOrder(userID int, orderID int, orderItems []model.OrderItems) (model.Order, error) {
	tx := c.db.MustBegin()

	o := model.Order{}
	rows, err := tx.NamedQuery(
		`UPDATE orders SET updated_at = now()
		WHERE user_id = :user_id AND id = :order_id RETURNING *`, map[string]interface{}{
			"user_id":  userID,
			"order_id": orderID,
		})
	if err != nil {
		return o, err
	}
	if rows.Next() {
		err := rows.StructScan(&o)
		if err != nil {
			tx.Rollback()
			return o, err
		}
	}

	rows.Close()

	// remove existing items from order
	_, err = tx.NamedExec(
		`UPDATE order_items SET deleted_at = now()
		WHERE order_id = :order_id AND deleted_at IS NULL`, map[string]interface{}{
			"order_id": orderID,
		})
	if err != nil {
		tx.Rollback()
		return o, err
	}

	for _, item := range orderItems {
		_, err = tx.NamedExec(
			`INSERT INTO order_items (order_id, coffee_id, quantity, created_at, updated_at) 
			VALUES (:order_id, :coffee_id, :quantity, now(), now())`, map[string]interface{}{
				"order_id":  o.ID,
				"coffee_id": item.Coffee.ID,
				"quantity":  item.Quantity,
			})
		if err != nil {
			tx.Rollback()
			return o, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return o, err
	}

	orders, err := c.GetOrders(userID, &orderID)
	if err != nil {
		return o, err
	}

	if len(orders) > 0 {
		return o, err
	}

	return orders[0], nil
}

// DeleteOrder deletes an existing order in the database
func (c *PostgresSQL) DeleteOrder(userID int, orderID int) error {
	tx := c.db.MustBegin()

	// remove existing items from order
	_, err := tx.NamedExec(
		`UPDATE order_items SET deleted_at = now()
		WHERE order_id = :order_id AND deleted_at IS NULL`, map[string]interface{}{
			"order_id": orderID,
		})
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(
		`UPDATE orders SET deleted_at = now()
		WHERE user_id = :user_id AND id = :order_id AND deleted_at IS NULL`, map[string]interface{}{
			"user_id":  userID,
			"order_id": orderID,
		})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CreateCoffee creates a new coffee
func (c *PostgresSQL) CreateCoffee(coffee model.Coffee) (model.Coffee, error) {
	m := model.Coffee{}

	rows, err := c.db.NamedQuery(
		`INSERT INTO coffees (name, teaser, description, price, image, created_at, updated_at) 
		VALUES(:name, :teaser, :description, :price, :image, now(), now()) 
		RETURNING id;`, map[string]interface{}{
			"name":        coffee.Name,
			"teaser":      coffee.Teaser,
			"description": coffee.Description,
			"price":       coffee.Price,
			"image":       coffee.Image,
		})
	if err != nil {
		return m, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&m)
		if err != nil {
			return m, err
		}
	}

	return m, nil
}

// CreateCoffeeIngredient creates a new coffee ingredient
func (c *PostgresSQL) CreateCoffeeIngredient(coffee model.Coffee, ingredient model.Ingredient) (model.CoffeeIngredients, error) {
	i := model.CoffeeIngredients{}

	rows, err := c.db.NamedQuery(
		`INSERT INTO coffee_ingredients (coffee_id, ingredient_id, quantity, unit,  created_at, updated_at) 
		VALUES(:coffee_id, :ingredient_id, :quantity, :unit, now(), now()) 
		RETURNING id;`, map[string]interface{}{
			"coffee_id":     coffee.ID,
			"ingredient_id": ingredient.ID,
			"quantity":      ingredient.Quantity,
			"unit":          ingredient.Unit,
		})
	if err != nil {
		return i, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&i)
		if err != nil {
			return i, err
		}
	}

	return i, nil
}
