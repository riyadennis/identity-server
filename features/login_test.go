package features

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/riyadennis/identity-server/internal/store"
	"golang.org/x/crypto/bcrypt"
)

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

func passwordWithFirstNameAndLastName(password, firstName, lastName string) error {
	enPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	userPassword = password
	user.Password = string(enPass)
	user.FirstName = firstName
	user.LastName = lastName
	return Idb.Insert(user)
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
