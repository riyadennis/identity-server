package store

import (
	"errors"
	// initialise mysql driver
	// initialise migration settings
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
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
		// TODO: add ping failure and valid config
		// {
		// 	name: "ping failure",
		// 	dbConn: &DBConnection{
		// 		User:     "test",
		// 		Password: "testPassword",
		// 		Host:     "localhost",
		// 		Port:     "3309",
		// 		Database: "test",
		// 	},
		// 	expectedErr: errPingFailed,
		// },
		// {
		// 	name:        "valid config",
		// 	dbConn:      NewENVConfig().DB,
		// 	expectedErr: nil,
		// },
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
