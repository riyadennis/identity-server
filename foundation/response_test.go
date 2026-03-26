package foundation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponse(t *testing.T) {
	scenarios := []struct {
		name      string
		status    int
		message   string
		errorCode string
	}{
		{
			name:      "bad request response",
			status:    http.StatusBadRequest,
			message:   "something went wrong",
			errorCode: "bad-request",
		},
		{
			name:      "unauthorised response",
			status:    http.StatusUnauthorized,
			message:   "not allowed",
			errorCode: UnAuthorised,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			r := NewResponse(sc.status, sc.message, sc.errorCode)
			assert.Equal(t, sc.status, r.Status)
			assert.Equal(t, sc.message, r.Message)
			assert.Equal(t, sc.errorCode, r.ErrorCode)
		})
	}
}

func TestJSONResponse(t *testing.T) {
	scenarios := []struct {
		name        string
		status      int
		message     string
		errorCode   string
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "ok response",
			status:      http.StatusOK,
			message:     "ok",
			errorCode:   "",
			wantStatus:  http.StatusOK,
			wantMessage: "ok",
		},
		{
			name:        "not found response",
			status:      http.StatusNotFound,
			message:     "resource not found",
			errorCode:   "not-found",
			wantStatus:  http.StatusNotFound,
			wantMessage: "resource not found",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := JSONResponse(rr, sc.status, sc.message, sc.errorCode)
			require.NoError(t, err)
			assert.Equal(t, sc.wantStatus, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			var got Response
			require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
			assert.Equal(t, sc.wantStatus, got.Status)
			assert.Equal(t, sc.wantMessage, got.Message)
		})
	}
}

func TestErrorResponse(t *testing.T) {
	scenarios := []struct {
		name        string
		status      int
		err         error
		errorCode   string
		wantMessage string
	}{
		{
			name:        "unauthorised",
			status:      http.StatusUnauthorized,
			err:         errors.New("not allowed"),
			errorCode:   UnAuthorised,
			wantMessage: "not allowed",
		},
		{
			name:        "internal server error",
			status:      http.StatusInternalServerError,
			err:         errors.New("database failure"),
			errorCode:   DatabaseError,
			wantMessage: "database failure",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			ErrorResponse(rr, sc.status, sc.err, sc.errorCode)
			assert.Equal(t, sc.status, rr.Code)

			var got Response
			require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
			assert.Equal(t, sc.status, got.Status)
			assert.Equal(t, sc.wantMessage, got.Message)
			assert.Equal(t, sc.errorCode, got.ErrorCode)
		})
	}
}

func TestResource(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	scenarios := []struct {
		name     string
		status   int
		resource payload
		wantName string
	}{
		{
			name:     "created resource",
			status:   http.StatusCreated,
			resource: payload{Name: "alice"},
			wantName: "alice",
		},
		{
			name:     "ok resource",
			status:   http.StatusOK,
			resource: payload{Name: "bob"},
			wantName: "bob",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := Resource(rr, sc.status, sc.resource)
			require.NoError(t, err)
			assert.Equal(t, sc.status, rr.Code)
			assert.Equal(t, "application/vnd.api+json; charset=utf-8", rr.Header().Get("Content-Type"))

			body, _ := io.ReadAll(rr.Body)
			var got payload
			require.NoError(t, json.Unmarshal(body, &got))
			assert.Equal(t, sc.wantName, got.Name)
		})
	}
}
