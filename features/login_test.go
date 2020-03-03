package features

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var (
	user *store.User
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

func password(arg1 string) error {
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

func iShouldGetErrorcode(arg1 string) error {
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

func aRegisteredUserWithEmailWithFirstNameAndLastName(email, fname, lname string) error {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/register", HOST),
		bytes.NewBuffer(registerInput(email, fname, lname)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	regResp, err := httpResponse(req)
	if err != nil {
		return err
	}
	password := strings.Split(regResp.Message, ":")
	user = &store.User{
		FirstName:        fname,
		LastName:         lname,
		Email:           email,
		Password:         strings.TrimSpace(password[1]),
	}
	return nil
}

func thatUserLogin() error {
	loginReq, err := http.NewRequest("POST",
		fmt.Sprintf("%s/login", HOST),
		bytes.NewBuffer(loginInput(user.Email, user.Password)))
	if err != nil {
		return err
	}
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err = httpResponse(loginReq)
	if err != nil {
		return err
	}
	logrus.Infof("%v", loginResp)
	return nil
}


func tokenNot(arg1 string) error {
	return godog.ErrPending
}
