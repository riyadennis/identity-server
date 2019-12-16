package handlers

import (
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	scenarios := []struct {
		name     string
		request  *http.Request
		response *Response
	}{
		{
			name: "no token",
			request: &http.Request{},
			response: &Response{
				Status:    http.StatusUnauthorized,
				Message:   "missing token",
				ErrorCode: UnAuthorised,
			},

		},
		{
			name: "invalid token",
			request: &http.Request{Header: http.Header{
				"Token": []string{
					"INVALID",
				},
			}},
			response: &Response{
				Status:    http.StatusUnauthorized,
				Message:   "invalid token",
				ErrorCode: UnAuthorised,
			},

		},
		{
			name: "invalid jwt token",
			request: &http.Request{Header: http.Header{
				"Token": []string{
					`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
						eyJleHAiOjE1NzY1MTIzMzZ9.
						1bBm4saeFtVRK8wMxuwApkLQUnbiGJYdpPrSQmua3Wo`,
				},
			}},
			response: &Response{
				Status:    http.StatusUnauthorized,
				Message:   "invalid token",
				ErrorCode: UnAuthorised,
			},

		},
	}
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			Home(rr, sc.request, nil)
			re := response(t, rr.Body)
			if !cmp.Equal(re, sc.response) {
				t.Errorf("expected %v, got %v", sc.response, re)
			}
		})
	}
}