Feature: Basic Functionality
  In order to ensure functionality defined in https://github.com/hashicorp-demoapp/product-api-go/issues/1
  Test the system

  Scenario: Get products
    Given the server is running
    When I make a "GET" request to "/coffees"
    Then A list of products should be returned