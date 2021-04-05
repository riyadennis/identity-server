package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	// initialise mysql driver
	_ "github.com/go-sql-driver/mysql"
	// initialise migration settings
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env_test")
	if err != nil {
		logrus.Fatalf("failed to open env file: %v", err)
	}
	cfg := store.NewENVConfig()
	conn, err = store.Connect(cfg.DB)
	if err != nil {
		logrus.Fatalf("failed to connect to db: %v", err)
	}

	err = store.Migrate(conn, cfg.DB.Database, cfg.BasePath)
	if err != nil {
		logrus.Fatalf("failed to run migration: %v", err)
	}

	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	scenarios := []struct {
		name             string
		req              *http.Request
		conn             *sql.DB
		expectedResponse *foundation.Response
	}{
		{
			name: "empty request",
			req:  nil,
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
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
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing email",
				foundation.ValidationFailed),
		},
		{
			name: "missing first name",
			req: registerPayLoad(t, &store.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     "joh@doe.com",
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing first name",
				foundation.ValidationFailed),
		},
		{
			name: "missing last name",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "",
				Email:     "joh@doe.com",
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing first name",
				foundation.ValidationFailed),
		},
		{
			name: "missing terms",
			req: registerPayLoad(t, &store.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "joh@doe.com",
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing terms",
				foundation.ValidationFailed),
		},
		{
			name: "invalid email",
			req: func() *http.Request {
				u := user(t)
				u.Email = "joh@dom"
				return registerPayLoad(t, u)
			}(),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"invalid email",
				foundation.ValidationFailed),
		},
		{
			name: "valid data empty connection",
			req:  registerPayLoad(t, user(t)),
			expectedResponse: foundation.NewResponse(http.StatusInternalServerError,
				"your generated password : GeneratedPassword",
				foundation.DatabaseError),
		},
		{
			name: "valid data empty connection",
			req:  registerPayLoad(t, user(t)),
			conn: nil,
			expectedResponse: foundation.NewResponse(http.StatusInternalServerError,
				"your generated password : GeneratedPassword",
				foundation.DatabaseError),
		},
		{
			name:             "success",
			req:              registerPayLoad(t, user(t)),
			conn:             conn,
			expectedResponse: foundation.NewResponse(http.StatusOK, "", ""),
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := &Handler{store.NewDB(sc.conn)}
			h.Register(w, sc.req, nil)
			resp := responseFromHTTP(t, w.Body)
			// TODO assert message also
			assert.Equal(t, sc.expectedResponse.ErrorCode, resp.ErrorCode)
			assert.Equal(t, sc.expectedResponse.Status, resp.Status)
		})
	}
}

func registerPayLoad(t *testing.T, u *store.User) *http.Request {
	jB, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(jB))
	req.Header.Set("content-type", "application/json")
	return req
}

func responseFromHTTP(t *testing.T, data io.Reader) *foundation.Response {
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
