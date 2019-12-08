package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	scenarios := []struct{
		name string
		request  *http.Request
		expectedStatus int
	}{
		{
			name: "invalid content type",
			request: &http.Request{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty content",
			request: request(t, `{hello: hi}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			request: request(t, `{
	"first_name": "John",
	"last_name": "Doe"
}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing terms",
			request: request(t, `{
	"first_name": "John",
	"last_name": "Doe",
	"email": "john@gmail.com"
}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "valid payload",
			request: request(t, `
{
	"first_name": "John",
	"last_name": "Doe",
	"email": "john@gmail.com",
	"terms": true
}`),
			expectedStatus: http.StatusOK,
		},
	}
	for _, sc := range scenarios{
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			Register(rr, sc.request, nil)
			if rr.Code != sc.expectedStatus{
				t.Errorf("want status code %d, got %d",
					sc.expectedStatus, rr.Code)
			}
		})
	}
}

func request(t *testing.T, content string) *http.Request{
	t.Helper()
	body := strings.NewReader(content)
	req, err:= http.NewRequest("POST",
		"/register", body)
	if err != nil{
		t.Error(err)
	}
	req.Header.Set("content-type", "application/json")
	return req
}