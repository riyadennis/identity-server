package main

import (
	"github.com/riyadennis/identity-server/business/store"
	"os"

	// initialise mysql driver
	_ "github.com/go-sql-driver/mysql"
	// initialise sqllite driver
	_ "github.com/mattn/go-sqlite3"
	// initialise migration settings
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/app/auth-api/handlers"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("failed to open env file: %v", err)
	}

	if os.Getenv("PORT") == "" {
		logrus.Fatal("invalid setting no port number")
	}

	if os.Getenv("ENV") == "" {
		logrus.Fatal("invalid setting no environment")
	}

	if os.Getenv("ENV") != "test" {
		file, err := os.Create("identity.log")
		if err != nil {
			logrus.Fatalf("failed to create log file: %v", err)
		}

		logrus.SetOutput(file)
	}
}

func main() {
	db, err := store.Connect()
	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = store.Migrate(db)
	if err != nil {
		logrus.Errorf("migration failed: %v", err)
	}

	err = handlers.NewServer(os.Getenv("PORT")).Run(db)
	logrus.Errorf("error running server: %v", err)
}
