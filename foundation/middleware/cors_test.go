package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	called := false
	allowedOrigins := []string{"http://example.com"}
	h := CORS(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		called = true
	}, allowedOrigins)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com") // Add Origin header
	rec := httptest.NewRecorder()

	h(rec, req, nil)
	assert.True(t, called, "next handler should be called")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Origin"), "http://example.com")
}
