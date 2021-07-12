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
    When I make a "POST" request to "/signout" where the request header is "Authorization" with the value "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjYxODc2NjYsInRva2VuX2lkIjoyLCJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6IlVzZXIxIn0.xprj5axiIs2NyrIVHJt_1G4eG3cDkHC1p84vhOsx4JI":
    Then the response status should be "OK"
