package store

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetUpMYSQL(t *testing.T) {
	// These tests focus on validation errors before attempting database connections
	// to ensure fast execution without network timeouts
	scenarios := []struct {
		name        string
		setupEnv    func()
		expectedErr error
	}{
		{
			name: "missing username",
			setupEnv: func() {
				os.Clearenv()
				_ = os.Setenv("MYSQL_PASSWORD", "password")
				_ = os.Setenv("MYSQL_HOST", "localhost")
				_ = os.Setenv("MYSQL_DATABASE", "identity")
				_ = os.Setenv("MYSQL_PORT", "3306")
			},
			expectedErr: errEmptyDBUserName,
		},
		{
			name: "missing password",
			setupEnv: func() {
				os.Clearenv()
				_ = os.Setenv("MYSQL_USERNAME", "root")
				_ = os.Setenv("MYSQL_HOST", "localhost")
				_ = os.Setenv("MYSQL_DATABASE", "identity")
				_ = os.Setenv("MYSQL_PORT", "3306")
			},
			expectedErr: errEmptyDBPassword,
		},
		{
			name: "missing host",
			setupEnv: func() {
				os.Clearenv()
				_ = os.Setenv("MYSQL_USERNAME", "root")
				_ = os.Setenv("MYSQL_PASSWORD", "password")
				_ = os.Setenv("MYSQL_DATABASE", "identity")
				_ = os.Setenv("MYSQL_PORT", "3306")
			},
			expectedErr: errEmptyDBHost,
		},
		{
			name: "missing database name",
			setupEnv: func() {
				os.Clearenv()
				_ = os.Setenv("MYSQL_USERNAME", "root")
				_ = os.Setenv("MYSQL_PASSWORD", "password")
				_ = os.Setenv("MYSQL_HOST", "localhost")
				_ = os.Setenv("MYSQL_PORT", "3306")
			},
			expectedErr: errEmptyDBName,
		},
		{
			name: "missing port",
			setupEnv: func() {
				os.Clearenv()
				_ = os.Setenv("MYSQL_USERNAME", "root")
				_ = os.Setenv("MYSQL_PASSWORD", "password")
				_ = os.Setenv("MYSQL_HOST", "localhost")
				_ = os.Setenv("MYSQL_DATABASE", "identity")
			},
			expectedErr: errEmptyDBPort,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			// Setup environment variables for this test case
			sc.setupEnv()

			// Execute
			store, auth, err := SetUpMYSQL(logrus.New())

			// Assert
			assert.Error(t, err, "Expected error for scenario: %s", sc.name)
			assert.ErrorIs(t, err, sc.expectedErr, "Expected specific error: %v", sc.expectedErr)
			assert.Nil(t, store, "Store should be nil on error")
			assert.Nil(t, auth, "Auth should be nil on error")
		})
	}

	// Cleanup environment
	os.Clearenv()
}

func TestSetUpMYSQL_AllMissingEnvVars(t *testing.T) {
	// Clear all environment variables
	os.Clearenv()

	logger := logrus.New()
	store, auth, err := SetUpMYSQL(logger)

	assert.Error(t, err)
	assert.Nil(t, store)
	assert.Nil(t, auth)
	assert.Contains(t, err.Error(), "mysql", "Error should mention mysql configuration issue")
}
