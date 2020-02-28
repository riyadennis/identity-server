Feature: user login
  In order to authenticate
  As a valid user
  I need to be able login

Scenario: user not registered
   Given email "john.doe@gmail.com"
   And password "password"
   When I login
   Then I should get error-code "invalid-request"
   And status code 400
   And message "invalid content"
