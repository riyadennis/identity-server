package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/foundation"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIdb struct {
	mock.Mock
}

func init() {
	Idb = &MockIdb{}
}

func (id *MockIdb) Insert(_ *store.User) error {
	return nil
}

func (id *MockIdb) Read(_ string) (*store.User, error) {
	return nil, nil
}

func (id *MockIdb) Authenticate(_, _ string) (bool, error) {
	return true, nil
}

func (id *MockIdb) Delete(_ string) (int64, error) {
	return 0, nil
}

func TestRegister(t *testing.T) {
	scenarios := []struct {
		name             string
		req              *http.Request
		expectedResponse *foundation.Response
	}{
		{
			name: "empty request",
			req:  nil,
			expectedResponse: foundation.NewResponse(400,
				"empty request",
				"invalid-request"),
		},
		{
			name: "missing email",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
			}),
			expectedResponse: foundation.NewResponse(400,
				"missing email",
				"invalid-user-data"),
		},
		{
			name: "missing first name",
			req: registerPayLoad(t, &store.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     "joh@doe.com",
			}),
			expectedResponse: foundation.NewResponse(400,
				"missing first name",
				"invalid-user-data"),
		},
		{
			name: "missing last name",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "",
				Email:     "joh@doe.com",
			}),
			expectedResponse: foundation.NewResponse(400,
				"missing first name",
				"invalid-user-data"),
		},
		{
			name: "missing terms",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "joh@doe.com",
			}),
			expectedResponse: foundation.NewResponse(400,
				"missing terms",
				"invalid-user-data"),
		},
		{
			name: "invalid email",
			req: func() *http.Request {
				u := user(t)
				u.Email = "joh@dom"
				return registerPayLoad(t, u)
			}(),
			expectedResponse: foundation.NewResponse(400,
				"invalid email",
				"invalid-user-data"),
		},
		{
			name: "valid data",
			req:  registerPayLoad(t, user(t)),
			expectedResponse: foundation.NewResponse(200,
				"your generated password : GeneratedPassword",
				""),
		},
	}
	w := httptest.NewRecorder()
	for _, sc := range scenarios {
		Register(w, sc.req, nil)
		resp := responseFromHttp(t, w.Body)
		// TODO assert message also
		assert.Equal(t, sc.expectedResponse.ErrorCode, resp.ErrorCode)
		assert.Equal(t, sc.expectedResponse.Status, resp.Status)
	}
}

func registerPayLoad(t *testing.T, u *store.User) *http.Request {
	jB, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("POST", "/register",
		bytes.NewReader(jB))
	req.Header.Set("content-type", "application/json")
	return req
}

func responseFromHttp(t *testing.T, data io.Reader) *foundation.Response {
	resp := &foundation.Response{}
	b, err := ioutil.ReadAll(data)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		t.Error(err)
	}
	return resp
}

func TestGeneratePassword(t *testing.T) {
	pass, err := business.GeneratePassword()
	if err != nil {
		t.Error(err)
	}
	if pass == "" {
		t.Error("empty password")
	}
}

func TestUserDataFromRequest(t *testing.T) {
	scenarios := []struct {
		name          string
		request       *http.Request
		expectedUser  *store.User
		expectedError string
	}{
		{
			name:          "empty request",
			request:       nil,
			expectedUser:  nil,
			expectedError: "empty request",
		},
		{
			name: "missing content type",
			request: &http.Request{
				Header: nil,
			},
			expectedUser:  nil,
			expectedError: "invalid content type",
		},
		{
			name: "missing body",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/register", nil)
				req.Header.Set("content-type", "application/json")
				return req
			}(),
			expectedUser:  nil,
			expectedError: "unexpected end of JSON input",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			u, err := userDataFromRequest(sc.request)
			assert.Equal(t, sc.expectedUser, u)
			if err != nil {
				assert.Equal(t, sc.expectedError, err.Error())
			}
		})
	}
}

func user(t *testing.T) *store.User {
	t.Helper()
	u := &store.User{
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "joh@doe.com",
		Password:         "testPassword",
		Company:          "testCompany",
		PostCode:         "E112QD",
		Terms:            true,
		RegistrationDate: "",
	}
	return u
}
