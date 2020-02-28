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
	userEmail    string
	userPassword string
	loginResp    *response
)

// response is the response we get back
// from rest call to login endpoint
type response struct {
	// Status is the http status code
	Status int `json:"status"`
	// Message is the message like
	// `welcome John Doe`
	Message string `json:"message"`
	// ErrorCode helps to debug issues
	// will be empty on success requests.
	ErrorCode string `json:"error-code"`
}

func email(arg1 string) error {
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
		bytes.NewBuffer(loginInput()))
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

// httpResponse will submit and http request using
// http client and then unmarshal the response into
// a struct.
func httpResponse(req *http.Request) (*response, error) {
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	loginResp = &response{}
	err = json.Unmarshal(body, loginResp)
	if err != nil {
		return nil, err
	}
	return loginResp, nil
}

// loginInput returns json request
// body for login endpoint in bytes.
func loginInput() []byte {
	return []byte(`{
	"email": "` + userEmail + `",
	"password": "` + userPassword + `"
}`)
}
