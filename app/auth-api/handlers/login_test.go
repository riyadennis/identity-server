package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				Status:    http.StatusBadRequest,
				Message:   "empty login data",
				ErrorCode: InvalidRequest,
			},
		},
		{
			name:    "missing email",
			request: request(t, "/login", `{}`),
			response: &Response{
				Status:    http.StatusBadRequest,
				Message:   "empty login data",
				ErrorCode: InvalidRequest,
			},
		},
		{
			name: "missing password",
			request: request(t, "/login", `{
	"email": "john4@gmail.com"
}`),
			response: &Response{
				Status:    http.StatusBadRequest,
				Message:   "empty login data",
				ErrorCode: InvalidRequest,
			},
		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			Login(rr, sc.request, nil)
			re := response(t, rr.Body)
			assert.Equal(t, sc.response, re)
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
