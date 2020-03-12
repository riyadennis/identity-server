Feature: user register
  In order to login
  As a valid user
  I need to be able register

  Scenario: new user register
    Given a not registered email "doe123@gmail.com"
    And firstName "John"
    And lastName "Doe"
    When that user register
    Then status code should be 200
    And errorCode ""
