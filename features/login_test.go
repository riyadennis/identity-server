package features

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ue  string
	up  string
	res *response
)

type response struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"error-code"`
}

func email(arg1 string) error {
	if arg1 == "" {
		return errors.New("empty email")
	}
	ue = arg1
	return nil
}

func password(arg1 string) error {
	if arg1 == "" {
		return errors.New("empty password")
	}
	up = arg1
	return nil
}

func iLogin() error {
	var jsonStr = []byte(`{
	"email": "johnmills@gmail.com",
	"password": "ePl6q3ARuOepz3c"
}`)
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/login", HOST),
		bytes.NewBuffer(jsonStr))
	req.Header.Set("authorization",
		"Basic amFjb2IyM0BnbWFpbC5jb206ZHJweEFqNUlvb0puS1dr'")
	client = &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	res = &response{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	return nil
}

func iShouldGetErrorcode(arg1 string) error {
	if res.ErrorCode != arg1 {
		return fmt.Errorf("expected error code %s, got %s",
			arg1, res.ErrorCode)
	}
	return nil
}

func statusCode(arg1 int) error {
	if res.Status != arg1 {
		return fmt.Errorf("expected status code %d, got %d",
			arg1, res.Status)
	}
	return nil
}

func message(arg1 string) error {
	if res.Message != arg1 {
		return fmt.Errorf("expected status code %s, got %s",
			arg1, res.Message)
	}
	return nil
}
