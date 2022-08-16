package data

import (
	"errors"
	"fmt"

	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/opentracing/opentracing-go/log"

	//"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection interface {
	IsConnected() (bool, error)
	GetCoffees(*int) (model.Coffees, error)
	GetIngredientsForCoffee(int) (model.Ingredients, error)
	CreateUser(string, string) (model.User, error)
	AuthUser(string, string) (model.User, error)
	CreateToken(int) (model.Token, error)
	GetToken(int, int) (model.Token, error)
	DeleteToken(int, int) error
	GetOrders(int, *int) (model.Orders, error)
	CreateOrder(int, []model.OrderItems) (model.Order, error)
	UpdateOrder(int, int, []model.OrderItems) (model.Order, error)
	DeleteOrder(int, int) error
	CreateCoffee(model.Coffee) (model.Coffee, error)
	UpsertCoffeeIngredient(model.Coffee, model.Ingredient) (model.CoffeeIngredient, error)
}

type PostgresSQL struct {
	telemetry *telemetry.Telemetry
	db        *sqlx.DB
}

// New creates a new connection to the database
func New(t *telemetry.Telemetry, connection string) (Connection, error) {
	t.AddMeasure("db.is_connected")
	t.AddMeasure("db.get_coffees")
	t.AddMeasure("db.get_ingredients_for_coffee")
	t.AddMeasure("db.create_user")
	t.AddMeasure("db.auth_user")
	t.AddMeasure("db.create_token")
	t.AddMeasure("db.get_token")
	t.AddMeasure("db.delete_token")
	t.AddMeasure("db.get_orders")
	t.AddMeasure("db.create_order")
	t.AddMeasure("db.update_order")
	t.AddMeasure("db.delete_order")
	t.AddMeasure("db.create_coffee")
	t.AddMeasure("db.upsert_coffee_ingredient")

	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &PostgresSQL{t, db}, nil
}

// IsConnected checks the connection to the database and returns an error if not connected
func (c *PostgresSQL) IsConnected() (bool, error) {
	done := c.telemetry.NewTiming("db.is_connected")
	defer done()

	err := c.db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetCoffees returns all coffees from the database
func (c *PostgresSQL) GetCoffees(coffeeid *int) (model.Coffees, error) {
	span := c.telemetry.Tracer.StartSpan("db.get_coffees")
	defer span.Finish()

	span.LogFields(log.String("db.system", "postgresql"))
	span.LogFields(
		log.String("event", "db.get_coffees"),
		log.String("value", "start"),
	)

	done := c.telemetry.NewTiming("db.get_coffees")
	defer done()

	cos := model.Coffees{}

	if coffeeid != nil {
		err := c.db.Select(&cos, "SELECT * FROM coffees WHERE id = $1", &coffeeid)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.db.Select(&cos, "SELECT * FROM coffees")
		if err != nil {
			return nil, err
		}
	}

	// fetch the ingredients for each coffee
	for n, cof := range cos {
		i := []model.CoffeeIngredient{}
		err := c.db.Select(&i, "SELECT ingredient_id FROM coffee_ingredients WHERE coffee_id=$1 AND quantity > 0", cof.ID)
		if err != nil {
			return nil, err
		}

		cos[n].Ingredients = i
	}

	span.LogFields(
		log.String("event", "db.get_coffees"),
		log.String("value", "finish"),
	)

	return cos, nil
}

// GetIngredientsForCoffee get the ingredients for the given coffeeid
func (c *PostgresSQL) GetIngredientsForCoffee(coffeeid int) (model.Ingredients, error) {
	done := c.telemetry.NewTiming("db.get_ingredients_for_coffee")
	defer done()

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
	done := c.telemetry.NewTiming("db.create_user")
	defer done()

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
	done := c.telemetry.NewTiming("db.auth_user")
	defer done()

	us := []model.User{}

	err := c.db.Select(&us,
		`SELECT id, username FROM users 
		WHERE username = $1 AND password = crypt($2, password);`,
		username, password,
	)
	if err != nil {
		return model.User{}, err
	}

	// If user does not exist
	if len(us) < 1 {
		return model.User{}, errors.New("User does not exist")
	}

	return us[0], nil
}

// CreateToken creates a new token
func (c *PostgresSQL) CreateToken(userID int) (model.Token, error) {
	done := c.telemetry.NewTiming("db.create_token")
	defer done()

	token := model.Token{}

	rows, err := c.db.NamedQuery(
		`INSERT INTO tokens (user_id, created_at) 
		VALUES(:user_id, now()) 
		RETURNING id;`, map[string]interface{}{
			"user_id": userID,
		})
	if err != nil {
		return token, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&token)
		if err != nil {
			return token, err
		}
	}

	return token, nil
}

// GetToken checks whether token exists
func (c *PostgresSQL) GetToken(tokenID int, userID int) (model.Token, error) {
	done := c.telemetry.NewTiming("db.get_token")
	defer done()

	token := []model.Token{}

	err := c.db.Select(&token,
		`SELECT id, user_id FROM tokens 
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;`,
		tokenID, userID,
	)
	if err != nil {
		return model.Token{}, err
	}

	if len(token) == 0 {
		return model.Token{}, fmt.Errorf("Invalid token")
	}

	return token[0], nil
}

// DeleteToken deletes an existing token in the database
func (c *PostgresSQL) DeleteToken(tokenID int, userID int) error {
	done := c.telemetry.NewTiming("db.delete_token")
	defer done()

	tx := c.db.MustBegin()

	_, err := tx.NamedExec(
		`UPDATE tokens SET deleted_at = now()
		WHERE id = :token_id AND user_id = :user_id AND deleted_at IS NULL`, map[string]interface{}{
			"token_id": tokenID,
			"user_id":  userID,
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

// GetOrders returns orders from the database
func (c *PostgresSQL) GetOrders(userID int, orderID *int) (model.Orders, error) {
	done := c.telemetry.NewTiming("db.get_orders")
	defer done()

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

				ing := []model.CoffeeIngredient{}
				err := c.db.Select(&ing, "SELECT ingredient_id FROM coffee_ingredients WHERE coffee_id=$1 AND quantity > 0", orders[n].Items[i].Coffee.ID)
				if err != nil {
					return nil, err
				}

				orders[n].Items[i].Coffee.Ingredients = ing
			}
		}
	}

	return orders, nil
}

// CreateOrder creates a new order in the database
func (c *PostgresSQL) CreateOrder(userID int, orderItems []model.OrderItems) (model.Order, error) {
	done := c.telemetry.NewTiming("db.create_order")
	defer done()

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
	done := c.telemetry.NewTiming("db.update_order")
	defer done()

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
	done := c.telemetry.NewTiming("db.delete_order")
	defer done()

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
	done := c.telemetry.NewTiming("db.create_coffee")
	defer done()

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

// UpsertCoffeeIngredient upserts a new coffee ingredient
func (c *PostgresSQL) UpsertCoffeeIngredient(coffee model.Coffee, ingredient model.Ingredient) (model.CoffeeIngredient, error) {
	done := c.telemetry.NewTiming("db.upsert_coffee_ingredient")
	defer done()

	i := model.CoffeeIngredient{}

	rows, err := c.db.NamedQuery(
		`INSERT INTO coffee_ingredients (coffee_id, ingredient_id, quantity, unit, created_at, updated_at) 
		VALUES(:coffee_id, :ingredient_id, :quantity, :unit, now(), now()) 
		ON CONFLICT ON CONSTRAINT unique_coffee_ingredient
		DO UPDATE SET quantity = :quantity, unit = :unit
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
