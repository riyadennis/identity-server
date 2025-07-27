package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func TestHandlerDelete(t *testing.T) {
	scenarios := []struct {
		name             string
		params           httprouter.Params
		expectedStatus   int
		expectedResponse *foundation.Response
	}{
		{
			name:           "no id",
			params:         nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   errInvalidID.Error(),
				ErrorCode: foundation.InvalidRequest,
			},
		},
		{
			name: "invalid id",
			params: httprouter.Params{
				{
					Key:   "id",
					Value: "INVALID",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &foundation.Response{
				Status:    http.StatusBadRequest,
				Message:   errDeleteFailed.Error(),
				ErrorCode: foundation.UserDoNotExist,
			},
		},
	}
	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			NewHandler(dbConn,
				&store.TokenConfig{
					Issuer:  "TEST",
					KeyPath: os.Getenv("KEY_PATH"),
				}, logger).
				Delete(w, nil, sc.params)

			assert.Equal(t, w.Code, sc.expectedStatus)

			resp := responseFromHTTP(t, w.Body)

			assert.Equal(t, sc.expectedResponse.ErrorCode, resp.ErrorCode)
			assert.Equal(t, sc.expectedResponse.Status, w.Code)
		})
	}
}
