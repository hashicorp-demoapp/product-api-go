Feature: Custom Coffee Functionality
  In order to ensure custom coffee functionality defined in https://github.com/hashicorp-demoapp/product-api-go/issues/14
  Test the system

  Scenario: Create a custom coffee
    Given the server is running
    When I make a "POST" request to "/coffees" where my userID is "1" with the following request body:
      """
      {"id":2,"name":"Latte", "teaser": "delicious custom coffee", "description": "best coffee in the world"}
      """
    Then a coffee should be returned
    And the response status should be "OK"

  Scenario: Create a custom coffee ingredient
    Given the server is running
    When I make a "POST" request to "/coffees/{id:[0-9]+}/ingredients" where my userID is "1" with the following request body:
      """
      {"coffee_id":2,"ingredient_id":1, "quantity": 50, "unit": "ml"}
      """
    Then a coffee ingredient should be returned
    And the response status should be "OK"