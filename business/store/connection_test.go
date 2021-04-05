package store

import (
	"database/sql"
	"errors"

	// initialise mysql driver
	_ "github.com/go-sql-driver/mysql"
	// initialise migration settings
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env_test")
	if err != nil {
		logrus.Fatalf("failed to open env file: %v", err)
	}
	cfg := NewENVConfig()
	conn, err = Connect(cfg.DB)
	if err != nil {
		logrus.Fatalf("failed to connect to db: %v", err)
	}

	err = Migrate(conn, cfg.DB.Database, cfg.BasePath)
	if err != nil {
		logrus.Fatalf("failed to run migration: %v", err)
	}

	os.Exit(m.Run())
}

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
			name: "empty db host",
			dbConn: &DBConnection{
				User:     "test",
				Password: "testPassword",
			},
			expectedErr: errEmptyDBHost,
		},
		{
			name: "empty db port",
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
			name: "ping failure",
			dbConn: &DBConnection{
				User:     "test",
				Password: "testPassword",
				Host:     "localhost",
				Port:     "3309",
				Database: "test",
			},
			expectedErr: errPingFailed,
		},
		{
			name:        "valid config",
			dbConn:      NewENVConfig().DB,
			expectedErr: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			_, err := Connect(sc.dbConn)
			if !errors.Is(err, sc.expectedErr) {
				t.Fatalf("unexpected error, wanted %v, got %v", sc.expectedErr, err)
			}
		})
	}
}
