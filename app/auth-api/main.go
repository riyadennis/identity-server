package main

import (
	"github.com/riyadennis/identity-server/app/auth-api/handlers"
	"github.com/riyadennis/identity-server/business/store"
	"os"

	// initialise mysql driver
	_ "github.com/go-sql-driver/mysql"
	// initialise sqllite driver
	_ "github.com/mattn/go-sqlite3"
	// initialise migration settings
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
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

	db, err := store.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	err = store.Migrate(db)
	if err != nil {
		logrus.Fatalf("migration failed: %v", err)
	}
}

func main() {
	handlers.Server(os.Getenv("PORT"))
}
