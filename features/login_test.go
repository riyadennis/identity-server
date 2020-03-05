package features

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/riyadennis/identity-server/internal/store"
)

const Password = "MUakRB5VndRu4U0"
var (
	user         *store.User
	userEmail    string
	userPassword string
	loginResp    *response
)

func aNotRegisteredEmail(arg1 string) error {
	if arg1 == "" {
		return errors.New("empty email")
	}
	userEmail = arg1
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
		bytes.NewBuffer(loginInput(userEmail, userPassword)))
	if err != nil {
		return err
	}
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
		Email:email,
	}
	return nil
}

func passwordFirstNameAndLastName(password, firstName, lastName string) error {
	enPass, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(enPass)
	user.FirstName = firstName
	user.LastName = lastName
	return Idb.Insert(user)
}

func thatUserLogin() error {
	loginReq, err := http.NewRequest("POST", fmt.Sprintf("%s/login", HOST),nil )
	if err != nil {
		return err
	}
	loginReq.Header.Set("Authorization",
		"Basic am9obi5kb2VAZ21haWwuY29tOk1VYWtSQjVWbmRSdTRVMA==")
	loginResp, err = httpResponse(loginReq)
	if err != nil {
		return err
	}
	return nil
}

func tokenNot(arg1 string) error {
	if loginResp.Token == arg1{
		return errors.New("invalid token")
	}
	return nil
}
