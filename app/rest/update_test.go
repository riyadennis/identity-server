package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/app/mocks"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
	"github.com/riyadennis/identity-server/foundation/middleware"
)

func TestHandlerUpdateUser(t *testing.T) {
	adminID := "admin-123"
	userID := "user-456"

	existingUser := &store.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Company:   "OldCo",
		PostCode:  "12345",
		CreatedBy: adminID,
		Active:    true,
	}

	updatedUser := &store.User{
		ID:        userID,
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "john@example.com",
		Company:   "NewCo",
		PostCode:  "67890",
		CreatedBy: adminID,
		Active:    true,
	}

	scenarios := []struct {
		name           string
		userID         string
		adminID        string
		body           interface{}
		mockStore      *mocks.Store
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "no user ID",
			userID:         "",
			adminID:        adminID,
			mockStore:      &mocks.Store{User: existingUser},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   foundation.InvalidRequest,
		},
		{
			name:           "user not found",
			userID:         userID,
			adminID:        adminID,
			mockStore:      &mocks.Store{User: nil},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   foundation.UserDoNotExist,
		},
		{
			name:    "not the creator",
			userID:  userID,
			adminID: "other-admin",
			mockStore: &mocks.Store{
				User: existingUser,
			},
			expectedStatus: http.StatusForbidden,
			expectedCode:   foundation.UnAuthorised,
		},
		{
			name:    "success",
			userID:  userID,
			adminID: adminID,
			body: map[string]string{
				"first_name": "Jane",
				"last_name":  "Smith",
				"company":    "NewCo",
				"post_code":  "67890",
			},
			mockStore: &mocks.Store{
				User: updatedUser,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			var bodyBytes []byte
			if sc.body != nil {
				bodyBytes, _ = json.Marshal(sc.body)
			} else {
				bodyBytes = []byte(`{"first_name":"Jane"}`)
			}

			w := httptest.NewRecorder()
			handler := NewHandler(sc.mockStore, &mocks.Authenticator{},
				&store.TokenConfig{}, logrus.New())

			r := httptest.NewRequest(http.MethodPut, "/admin/update/"+sc.userID,
				bytes.NewReader(bodyBytes))

			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("userID", sc.userID)

			claims := &jwt.RegisteredClaims{Subject: sc.adminID}
			ctx := context.WithValue(r.Context(), chi.RouteCtxKey, routeContext)
			ctx = context.WithValue(ctx, middleware.UserClaimsKey, claims)
			r = r.WithContext(ctx)

			handler.UpdateUser(w, r)

			assert.Equal(t, sc.expectedStatus, w.Code)

			if sc.expectedCode != "" {
				resp := responseFromHTTP(t, w.Body)
				if resp != nil {
					assert.Equal(t, sc.expectedCode, resp.ErrorCode)
				}
			}
		})
	}
}
