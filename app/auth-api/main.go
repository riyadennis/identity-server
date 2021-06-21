package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	// initialise mysql driver
	_ "github.com/go-sql-driver/mysql"
	// initialise migration settings
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/riyadennis/identity-server/app/auth-api/handlers"
	"github.com/riyadennis/identity-server/business/store"
)

func main() {
	useEnvFile := os.Args[1]

	logger := log.New(os.Stdout, "IDENTITY: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	fmt.Printf("args: %s", useEnvFile)

	if useEnvFile == "true" {
		err := godotenv.Load()
		if err != nil {
			logger.Fatalf("failed to open env file: %v", err)
		}
	}

	logger.Printf("env file switch value %v", useEnvFile)
	cfg := store.NewENVConfig()

	db, err := store.Connect(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = store.Migrate(db, cfg.DB.Database, cfg.BasePath)
	if err != nil {
		logger.Panicf("migration failed: %v", err)
	}

	err = handlers.NewServer(os.Getenv("PORT")).Run(db, cfg.Token, logger)
	if err != nil {
		logger.Panicf("error running server: %v", err)
	}
}
