package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/jsonapi"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

var (
	dbConn     *store.DB
	testLogger *log.Logger
	testEmail  string
)

func TestRegister(t *testing.T) {
	scenarios := []struct {
		name             string
		req              *http.Request
		conn             *sql.DB
		expectedResponse *foundation.Response
	}{
		{
			name: "missing email",
			req: registerPayLoad(t, &store.UserRequest{
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
			req: registerPayLoad(t, &store.UserRequest{
				FirstName: "",
				LastName:  "Doe",
				Email:     testEmail,
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing first name",
				foundation.ValidationFailed),
		},
		{
			name: "missing last name",
			req: registerPayLoad(t, &store.UserRequest{
				FirstName: "John",
				LastName:  "",
				Email:     testEmail,
			}),
			expectedResponse: foundation.NewResponse(http.StatusBadRequest,
				"missing first name",
				foundation.ValidationFailed),
		},
		{
			name: "missing terms",
			req: registerPayLoad(t, &store.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     testEmail,
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
	}
	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			db := store.NewDB(sc.conn)
			h := NewHandler(db,
				&store.TokenConfig{
					Issuer:  "TEST",
					KeyPath: os.Getenv("KEY_PATH"),
				}, logger)
			h.Register(w, sc.req, nil)
			assert.Equal(t, sc.expectedResponse.Status, w.Code)
			resp := responseFromHTTP(t, w.Body)
			// TODO assert message also
			assert.Equal(t, sc.expectedResponse.ErrorCode, resp.ErrorCode)
			assert.Equal(t, sc.expectedResponse.Status, w.Code)
		})
	}
}

func registerPayLoad(t *testing.T, u *store.UserRequest) *http.Request {
	var buff bytes.Buffer

	err := jsonapi.MarshalPayload(&buff, u)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/register", strings.NewReader(buff.String()))
	req.Header.Set("content-type", "application/json")
	return req
}

func responseFromHTTP(t *testing.T, data io.Reader) *foundation.Response {
	t.Helper()

	resp := &foundation.Response{}
	b, err := io.ReadAll(data)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		t.Error(err)
	}

	if resp.Message != "" {
		return resp
	}

	return nil
}

func TestGeneratePassword(t *testing.T) {
	t.Helper()

	pass, err := business.GeneratePassword()
	if err != nil {
		t.Error(err)
	}
	if pass == "" {
		t.Error("empty password")
	}
}

func user(t *testing.T) *store.UserRequest {
	t.Helper()

	u := &store.UserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  "testPassword",
		Company:   "testCompany",
		PostCode:  "E112QD",
		Terms:     true,
	}
	return u
}
