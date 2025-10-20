package foundation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/stretchr/testify/assert"
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
}
