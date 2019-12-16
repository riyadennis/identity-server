package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogin(t *testing.T) {
	scenarios := []struct {
		name     string
		request  *http.Request
		response *Response
	}{
		{
			name:    "empty request",
			request: &http.Request{},
			response: &Response{
				Status:    400,
				Message:   "invalid content",
				ErrorCode: InvalidRequest,
			},
		},
		{
			name:    "missing email",
			request: request(t, `{}`),
			response: &Response{
				Status:    400,
				Message:   "email missing",
				ErrorCode: EmailMissing,
			},
		},
		{
			name: "missing password",
			request: request(t, `{
	"email": "john4@gmail.com"
}`),
			response: &Response{
				Status:    400,
				Message:   "password missing",
				ErrorCode: PassWordError,
			},
		},
		{
			name: "invalid password",
			request: request(t, `{
	"email": "john4@gmail.com",
	"password": "invalid"
}`),
			response: &Response{
				Status:    400,
				Message:   "sql: no rows in result set",
				ErrorCode: InvalidRequest,
			},
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			Login(rr, sc.request, nil)
			re := response(t, rr.Body)
			if !cmp.Equal(re, sc.response) {
				t.Errorf("expected %v, got %v", sc.response, re)
			}
		})
	}
}

func response(t *testing.T, body *bytes.Buffer) *Response {
	t.Helper()
	var re *Response
	if body != nil {
		data, err := ioutil.ReadAll(body)
		if err != nil {
			t.Error(err.Error())
		}
		re = &Response{}
		err = json.Unmarshal(data, re)
		if err != nil {
			t.Error(err.Error())
		}
	}
	return re
}
