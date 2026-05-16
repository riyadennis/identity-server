package foundation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestBody(t *testing.T) {
	scenarios := []struct {
		name          string
		request       *http.Request
		expectedUser  *store.User
		expectedError string
	}{
		{
			name: "missing content type",
			request: &http.Request{
				Header: nil,
			},
			expectedUser:  nil,
			expectedError: errMissingContentType.Error(),
		},
		{
			name: "missing body",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/register", nil)
				req.Header.Set("content-type", "application/json")
				return req
			}(),
			expectedUser:  nil,
			expectedError: "EOF",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := RequestBody(sc.request, sc.expectedUser)
			if err != nil {
				assert.Equal(t, sc.expectedError, err.Error())
			}
		})
	}

	t.Run("nil request", func(t *testing.T) {
		err := RequestBody(nil, &store.User{})
		assert.Equal(t, errEmptyRequest, err)
	})

	t.Run("success", func(t *testing.T) {
		body := `{"first_name":"John","last_name":"Doe","email":"john@test.com"}`
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("content-type", "application/json")

		var user store.User
		err := RequestBody(req, &user)
		require.NoError(t, err)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john@test.com", user.Email)
	})
}
