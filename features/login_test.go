package features

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"net/http"

	"github.com/riyadennis/identity-server/internal/store"
)

var (
	user         *store.User
	userEmail    string
	userPassword string
	loginResp    *response
)

func LoginFeatureContext(s *godog.Suite) {
	s.BeforeScenario(beforeScenario)
	s.Step(`^a not registered email "([^"]*)"$`, aNotRegisteredEmail)
	s.Step(`^password "([^"]*)"$`, aNotRegisteredPassword)
	s.Step(`^I login$`, iLogin)
	s.Step(`^I should get error-code "([^"]*)"$`, iShouldGetErrorCode)
	s.Step(`status code (\d+)$`, statusCode)
	s.Step(`^message "([^"]*)"$`, message)

	s.Step(`^a registered user with email "([^"]*)"$`, aRegisteredUserWithEmail)
	s.Step(`^password "([^"]*)" with firstName "([^"]*)" and lastName "([^"]*)""$`, passwordWithFirstNameAndLastName)

	s.Step(`^that user login$`, thatUserLogin)
	s.Step(`^status code should be (\d+)$`, statusCode)

	s.AfterScenario(afterScenario)
}

func aNotRegisteredEmail(email string) error {
	if email == "" {
		return errors.New("empty email")
	}
	userEmail = email
	return nil
}

func aNotRegisteredPassword(arg1 string) error {
	if arg1 == "" {
		return errors.New("empty password")
	}
	userPassword = arg1
	return nil
}

func iLogin() error {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/login", HOST),
		nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", getBase64())
	loginResp, err = httpResponse(req)
	if err != nil {
		return err
	}
	return nil
}

func iShouldGetErrorCode(arg1 string) error {
	if loginResp.ErrorCode != arg1 {
		return fmt.Errorf("expected error code %s, got %s",
			arg1, loginResp.ErrorCode)
	}
	return nil
}

func statusCode(arg1 int) error {
	if loginResp.Status != arg1 {
		return fmt.Errorf("expected status code %d, got %d",
			arg1, loginResp.Status)
	}
	return nil
}

func message(arg1 string) error {
	if loginResp.Message != arg1 {
		return fmt.Errorf("expected status code %s, got %s",
			arg1, loginResp.Message)
	}
	return nil
}

func aRegisteredUserWithEmail(email string) error {
	user = &store.User{
		Email: email,
	}
	userEmail = email
	return nil
}

func passwordWithFirstNameAndLastName(firstName, lastName string) error {
	user.FirstName = firstName
	user.LastName = lastName
	return thatUserRegister()
}

func thatUserLogin() error {
	loginReq, err := http.NewRequest("POST", fmt.Sprintf("%s/login", HOST), nil)
	if err != nil {
		return err
	}
	loginReq.Header.Set("Authorization", getBase64())
	loginResp, err = httpResponse(loginReq)
	if err != nil {
		return err
	}
	return nil
}

func getBase64() string {
	str := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", userEmail, userPassword)))
	return fmt.Sprintf("Basic %s", str)
}
