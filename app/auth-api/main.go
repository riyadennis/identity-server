package main

import (
	"log"
	"os"

	// initialise mysql driver
	_ "github.com/go-sql-driver/mysql"
	// initialise migration settings
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/app/auth-api/handlers"
	"github.com/riyadennis/identity-server/business/store"
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
	cfg := store.NewENVConfig()

	db, err := store.Connect(cfg.DB)
	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = store.Migrate(db, cfg.DB.Database, cfg.BasePath)
	if err != nil {
		logrus.Errorf("migration failed: %v", err)
	}

	logger := log.New(os.Stdout, "IDENTITY : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	err = handlers.NewServer(os.Getenv("PORT")).Run(db, logger)
	if err != nil {
		logrus.Errorf("error running server: %v", err)
	}
}
