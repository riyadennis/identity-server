package store

import (
	"database/sql"
	"errors"
	"testing"
)

func TestMigrate(t *testing.T) {
	scenarios := []struct {
		name        string
		dbConn      *sql.DB
		dbName      string
		basePath    string
		expectedErr error
	}{
		{
			name:        "empty database connection",
			dbConn:      nil,
			expectedErr: errEmptyDBConnection,
		},
		{
			name:        "empty database name",
			dbConn:      conn,
			dbName:      "",
			expectedErr: errEmptyDatabaseName,
		},
		{
			name:        "empty base path",
			dbConn:      conn,
			dbName:      "test",
			expectedErr: errMigrationInitialisation,
		},
		{
			name:        "success",
			dbConn:      conn,
			dbName:      "test",
			basePath:    "../../",
			expectedErr: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := Migrate(sc.dbConn, sc.dbName, sc.basePath)
			if !errors.Is(err, sc.expectedErr) {
				t.Fatalf("unexpected error, wanted %v, got %v", sc.expectedErr, err)
			}
		})
	}
}
