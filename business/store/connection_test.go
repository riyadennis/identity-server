package store

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	// initialise mysql driver
	// initialise migration settings
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	scenarios := []struct {
		name        string
		dbConn      *DBConnection
		expectedErr error
	}{
		{
			name:        "empty connection",
			dbConn:      nil,
			expectedErr: errInvalidDBConfig,
		},
		{
			name:        "empty user name",
			dbConn:      &DBConnection{},
			expectedErr: errEmptyDBUserName,
		},
		{
			name: "empty password",
			dbConn: &DBConnection{
				User: "test",
			},
			expectedErr: errEmptyDBPassword,
		},
		{
			name: "empty auth host",
			dbConn: &DBConnection{
				User:     "test",
				Password: "testPassword",
			},
			expectedErr: errEmptyDBHost,
		},
		{
			name: "empty auth port",
			dbConn: &DBConnection{
				User:     "test",
				Password: "testPassword",
				Host:     "localhost",
			},
			expectedErr: errEmptyDBPort,
		},
		{
			name: "empty database name",
			dbConn: &DBConnection{
				User:     "test",
				Password: "testPassword",
				Host:     "localhost",
				Port:     "3309",
			},
			expectedErr: errEmptyDBName,
		},
		{
			name: "valid config",
			dbConn: func() *DBConnection {
				_ = os.Setenv("MYSQL_USERNAME", "root")
				_ = os.Setenv("MYSQL_PASSWORD", "root")
				_ = os.Setenv("MYSQL_HOST", "localhost")
				_ = os.Setenv("MYSQL_DATABASE", "identity")
				_ = os.Setenv("MYSQL_PORT", "80")
				return NewENVConfig().DB
			}(),
			expectedErr: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			_, err := ConnectMYSQL(sc.dbConn)
			if !errors.Is(err, sc.expectedErr) {
				t.Fatalf("unexpected error, wanted %v, got %v", sc.expectedErr, err)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(os.Stderr)

	t.Run("invalid key bytes", func(t *testing.T) {
		_, err := GenerateToken(logger, []byte("not valid pem"), &jwt.RegisteredClaims{})
		assert.Error(t, err)
	})

	t.Run("valid RSA key generates token", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)

		privPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})

		expiry := time.Now().Add(1 * time.Hour)
		token, err := GenerateToken(logger, privPEM, &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    "test-issuer",
			Subject:   "user-123",
		})
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.AccessToken)
		assert.Equal(t, 200, token.Status)
		assert.Equal(t, "Bearer", token.TokenType)
	})
}
