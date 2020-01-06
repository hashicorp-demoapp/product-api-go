CREATE TABLE coffees (id serial PRIMARY KEY, name VARCHAR (255) NOT NULL, price NUMERIC(5, 2) NOT NULL, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL, deleted_at TIMESTAMP);
CREATE TABLE ingredients (id serial PRIMARY KEY, name VARCHAR (255) NOT NULL, quantity VARCHAR (50) NOT NULL, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL, deleted_at TIMESTAMP);
CREATE TABLE coffee_ingredients (id serial PRIMARY KEY, coffee_id int NOT NULL, ingredient_id int NOT NULL, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL, deleted_at TIMESTAMP);

INSERT INTO ingredients (id, name, quantity, created_at, updated_at) VALUES (1, 'Double shot espresso', '20ml', CURRENT_DATE, CURRENT_DATE);
INSERT INTO ingredients (id, name, quantity, created_at, updated_at) VALUES (2, 'Semi skimmed Milk', '500ml', CURRENT_DATE, CURRENT_DATE);
INSERT INTO ingredients (id, name, quantity, created_at, updated_at) VALUES (3, 'Hot Water', '500ml', CURRENT_DATE, CURRENT_DATE);


INSERT INTO coffees (name, price, created_at, updated_at) VALUES ('Latte (Medium)', 200.00, CURRENT_DATE, CURRENT_DATE);
INSERT INTO coffee_ingredients (coffee_id, Ingredient_id, created_at, updated_at) VALUES (1,1, CURRENT_DATE, CURRENT_DATE);
INSERT INTO coffee_ingredients (coffee_id, Ingredient_id, created_at, updated_at) VALUES (1,2, CURRENT_DATE, CURRENT_DATE);


INSERT INTO coffees (name, price, created_at, updated_at) VALUES ('Americano (Medium)', 150.00, CURRENT_DATE, CURRENT_DATE);
INSERT INTO coffee_ingredients (coffee_id, Ingredient_id, created_at, updated_at) VALUES (2,1, CURRENT_DATE, CURRENT_DATE);
INSERT INTO coffee_ingredients (coffee_id, Ingredient_id, created_at, updated_at) VALUES (2,3, CURRENT_DATE, CURRENT_DATE);