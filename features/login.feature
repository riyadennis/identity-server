Feature: user login
  In order to authenticate
  As a valid user
  I need to be able login

Scenario: user not registered
   Given a not registered email "john.doe@gmail.com"
   And password "password"
   When I login
   Then I should get error-code "invalid-request"
   And status code 400
   And message "empty login data"

Scenario: user registered login with invalid password
   Given a registered user with email "john.doe@gmail.com"
   And not matching password "INVALID" with firstName "John" and lastName "Doe""
   When that user login
   Then status code should be 400
   And message "crypto/bcrypt: hashedPassword is not the hash of the given password"

