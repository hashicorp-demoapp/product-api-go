Feature: Basic Order Functionality
  In order to ensure functionality defined in https://github.com/hashicorp-demoapp/product-api-go/issues/2
  Test the system

  Scenario: View all user orders
    Given the server is running
    When I make a "GET" request to "/orders" where my userID is "1"
    Then a list of orders should be returned
    And the response status should be "OK"

  Scenario: View a specific order
    Given the server is running
    When I make a "GET" request to "/orders/{id:[0-9]+}" where "id" is "1" and my userID is "1"
    Then an order should be returned
    And the response status should be "OK"

  Scenario: Create an order
    Given the server is running
    When I make a "POST" request to "/orders" where my userID is "1" with the following request body:
      """
      [{"coffee":{"id":1,"name":"Latte"},"quantity":2},{"coffee":{"id":2,"name":"Americano"},"quantity":3}]
      """
    Then an order should be returned
    And the response status should be "OK"

  Scenario: Update an order
    Given the server is running
    When I make a "PUT" request to "/orders/{id:[0-9]+}" where "id" is "1" and my userID is "1", with the following request body:
      """
      [{"coffee":{"id":1,"name":"Latte"},"quantity":2},{"coffee":{"id":2,"name":"Americano"},"quantity":3}]
      """
    Then an order should be returned
    And the response status should be "OK"

  Scenario: Delete an order
    Given the server is running
    When I make a "DELETE" request to "/orders/{id:[0-9]+}" where "id" is "1" and my userID is "1"
    Then the response status should be "OK"