Feature: Basic User Functionality
  In order to ensure functionality defined in https://github.com/hashicorp-demoapp/product-api-go/issues/2
  Test the system

  Scenario: Create new user
    Given the server is running
    When I make a "POST" request to "/signup" with the following request body:
      """
      {"username": "User1", "password": "testPassword"}
      """
    Then an AuthResponse should be returned
    And the response status should be "OK"

  Scenario: Sign in user
    Given the server is running
    When I make a "POST" request to "/signin" with the following request body:
      """
      {"username": "User1", "password": "testPassword"}
      """
    Then an AuthResponse should be returned
    And the response status should be "OK"
