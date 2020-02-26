package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
)

func init() {
	err := sqlite.Setup("/var/tmp/identityTest.db")
	if err != nil {
		panic(err)
	}
	db, err := sqlite.ConnectDB("/var/tmp/identityTest.db")
	if err != nil {
		panic(err)
	}
	Idb,  err = sqlite.PrepareDB(db)
	if err != nil {
		panic(err)
	}
}

func TestRegister(t *testing.T) {
	scenarios := []struct {
		name           string
		request        *http.Request
		expectedStatus int
	}{
		{
			name:           "invalid content type",
			request:        &http.Request{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty content",
			request:        request(t,"/register", `{hello: hi}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			request: request(t, "/register", `{
	"first_name": "John",
	"last_name": "Doe"
}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing terms",
			request: request(t, "/register", `{
	"first_name": "John",
	"last_name": "Doe",
	"email": "john@gmail.com"
}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing terms",
			request: request(t, "/register", `{
	"first_name": "John",
	"last_name": "Doe",
	"email": "@gml.com"
}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "valid payload",
			request: request(t,"/register",  `
{
	"first_name": "John",
	"last_name": "Doe",
	"email": "john@gmail.com",
	"terms": true
}`),
			expectedStatus: http.StatusOK,
		},
		{
			name: "duplicate payload",
			request: request(t, "/register", `
{
	"first_name": "John",
	"last_name": "Doe",
	"email": "john@gmail.com",
	"terms": true
}`),
			expectedStatus: http.StatusBadRequest,
		},
	}
	defer func() {
		err := os.Remove("/var/tmp/identityTest.db")
		if err != nil {
			panic(err)
		}
	}()
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			Register(rr, sc.request, nil)
			if rr.Code != sc.expectedStatus {
				t.Errorf("want Status code %d, got %d",
					sc.expectedStatus, rr.Code)
			}
		})
	}

}

func TestGeneratePassword(t *testing.T) {
	pass, err := generatePassword()
	if err != nil {
		t.Error(err)
	}
	if pass == "" {
		t.Error("empty password")
	}
}

func request(t *testing.T, endpoint, content string) *http.Request {
	t.Helper()
	body := strings.NewReader(content)
	req, err := http.NewRequest("POST",
		endpoint, body)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("content-type", "application/json")
	return req
}
