Feature: accounts
  In order to use accounts api
  As an API user
  I need to be able to manage accounts

  Scenario: should get empty accounts list if no accounts are created
    When I list accounts
    Then the response code should be 200
    And the count of accounts should be 0

  Scenario: should create an account
    When I create account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    Then the response code should be 201
    And the created account id should be "73c4ee80-e60e-11e9-a044-acde48001122"

  Scenario: should create and find account
    When I create account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    And I fetch account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    Then the response code should be 200
    And the fetched account id should be "73c4ee80-e60e-11e9-a044-acde48001122"

  Scenario: should delete created account
    When I create account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    And I delete account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    Then the response code should be 204

  Scenario: should not fetch deleted account
    When I create account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    And I delete account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    And I fetch account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    Then the response code should be 404

  Scenario: should not list deleted account
    When I create account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    And I delete account with id "73c4ee80-e60e-11e9-a044-acde48001122"
    And I list accounts
    Then the response code should be 200
    And the count of accounts should be 0

  Scenario: should list the correct amount accounts per page
    When I create 10 accounts
    And I list 5 accounts per page
    Then the response code should be 200
    And the count of accounts should be 5

  Scenario: should list the correct amount accounts per page
    When I create 10 accounts
    And I list 5 accounts per page
    Then the response code should be 200
    And the count of accounts should be 5
    And the current page should be 1
    And the page should not be the last

  Scenario: should list the correct amount accounts per page
    When I create 8 accounts
    And I list 5 accounts per page in page 2
    Then the response code should be 200
    And the count of accounts should be 3
    And the current page should be 1
    And the page should be the last