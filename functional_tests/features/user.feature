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

  Scenario: Sign out user
    Given the server is running
    When I make a "POST" request to "/signout" where the request header is "Authorization" with the value "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl9pZCI6MiwidXNlcl9pZCI6MSwidXNlcm5hbWUiOiJVc2VyMSJ9.D0DltVBYHrFwL4GJO3x1K0ZieFhmwAHcVLIMIzDQ-Ek":
    Then the response status should be "OK"
