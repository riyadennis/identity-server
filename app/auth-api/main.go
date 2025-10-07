package main

import (
	"log"
	"os"
	// initialise mysql driver
	// initialise migration settings

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"

	"github.com/riyadennis/identity-server/app/auth-api/handlers"
	"github.com/riyadennis/identity-server/business/store"
)

func main() {
	logger := log.New(os.Stdout, "IDENTITY: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("failed to open env file: %v", err)
	}

	cfg := store.NewENVConfig()

	db, err := store.Connect(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("database ping failed: %v", err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = store.Migrate(db, cfg.DB.Database, cfg.DB.MigrationPath)
	if err != nil {
		logger.Panicf("migration failed: %v", err)
	}

	err = handlers.NewServer(os.Getenv("PORT")).Run(db, cfg.Token, logger)
	if err != nil {
		logger.Panicf("error running server: %v", err)
	}
}
