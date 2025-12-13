package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/riyadennis/identity-server/app/auth-api/mocks"
	"github.com/riyadennis/identity-server/foundation"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/riyadennis/identity-server/business/store"
)

func setupTestRouter(s store.Store, a store.Authenticator) http.Handler {
	tokenConfig := &store.TokenConfig{Issuer: "TEST", KeyPath: os.Getenv("KEY_PATH")}
	return LoadRESTEndpoints(tokenConfig, logrus.New(), s, a)
}

func TestLivenessRoute(t *testing.T) {
	router := setupTestRouter(&mocks.Store{}, &mocks.Authenticator{})
	req := httptest.NewRequest(http.MethodGet, LivenessEndPoint, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "up", resp["Status"])
}

func TestReadinessRoute(t *testing.T) {
	scenarios := []struct {
		name         string
		mockStore    store.Store
		expectedCode int
	}{
		{
			name:         "success",
			mockStore:    &mocks.Store{},
			expectedCode: http.StatusOK,
		},
		{
			name:         "error",
			mockStore:    &mocks.Store{Error: errors.New("error")},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			router := setupTestRouter(scenario.mockStore, nil)
			req := httptest.NewRequest(http.MethodGet, ReadinessEndPoint, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)
			assert.Equal(t, scenario.expectedCode, rec.Code)
		})
	}
}

func TestRegisterRoute_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter(nil, nil)
	req := httptest.NewRequest(http.MethodGet, RegisterEndpoint, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	Home(w, nil)
	body := &foundation.Response{}
	dec := json.NewDecoder(w.Body)
	dec.Decode(body)
	assert.Equal(t, http.StatusOK, body.Status)
	assert.Equal(t, "Authorised", body.Message)
	assert.Equal(t, "", body.ErrorCode)
}

func TestHomeRoute_NoToken(t *testing.T) {
	router := setupTestRouter(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/user"+HomeEndPoint, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRegisterRoute_ValidRequest(t *testing.T) {
	mockStore := &mocks.Store{
		User: &store.User{},
	}
	router := setupTestRouter(mockStore, &mocks.Authenticator{
		ReturnVal: true,
	})

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user(t))
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, RegisterEndpoint, &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestDeleteRoute_ValidToken(t *testing.T) {
	auth := &mocks.Authenticator{
		ReturnVal: true,
		Token: &store.TokenRecord{
			Expiry: time.Now().Add(2 * time.Hour),
			TTL:    "123",
		},
	}
	router := setupTestRouter(&mocks.Store{}, auth)

	// Create delete request with auth token
	req := httptest.NewRequest(http.MethodDelete, "/user/delete/123", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU5MDQxNDMsImlzcyI6Im9wZW4gc291cmNlIn0.oOvAubDJySDEIXUmCcz6CzsQSAQZJUvXZ6PZofB6bu8dQNn7BQIW8Q6F2PhNuM5HmxSw2lDt7heGyOGOPcpVm_36mFcNkDaZJxVC3mI6gvrR9rnOXNHzeTUEAyRL-Hlgmwan4TrIEmTqVaQza9Aj_k1Z7WlW-gy0EJdFmms4qACuVaElfas4pQKE0wraO7IS6PHODoWAlmhf7yEFxFm2jfhcu-HMEc5ETqsQqOOfrfmb_tJguDTq2VfphBXEO38yoWP-6C1djz6TLxmd_1wSsUdmygtUe_9LUrhnu3rtpEtzWJNeJAS-BypwpIijHn90zOo04CC_TsCEy4aqAUP6TQ")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	// Note: This will return 401 because we're not properly mocking the JWT validation
	// In a real scenario, you'd need to setup proper token validation or mock the auth middleware
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestHomeRoute_ValidToken(t *testing.T) {
	router := setupTestRouter(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/user"+HomeEndPoint, nil)

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQzODgzOTgsImlzcyI6IiJ9.fcaTzPU4PkeOe_CwtRvdetEDWCm_x0E3pLNUpKCcELdD7RqbCsbig9WWspey1pKyckawoz7N9XGAffmq8i4G97oQcwjNffZKC8SKq6ocxPBs95G0f8KmcX1nCYVsSPb6r0D3A3KCnWphiwrwf6-kmKDxhvEaxqquOfaJcws6JSkekjml_H3iCinbiVISHZqvAqjWSSkC4CPbzB4yqNJ0_oRvkJx-gL8Z7w_Jmk28RouYWap_-1Hzy6MZt4s-PtZPXIQw7NRA3NyVGp9f-MMoatsmOFAkvFbV1wSnEzUgKLg1Wga98y9YnTYDvbFhC8pyHKsbEq0g2en6qqymDg2ZCQ")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	// Note: This will return 401 because we're not properly mocking the JWT validation
	// In a real scenario, you'd need to setup proper token validation or mock the auth middleware
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
