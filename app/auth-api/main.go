package main

import (
	"os"
	// initialise mysql driver
	// initialise migration settings

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/app/auth-api/server"
	"github.com/riyadennis/identity-server/business/store"
)

func main() {
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
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

	s := server.NewServer(os.Getenv("PORT"))
	err = s.Run(db, cfg.Token, logger)

	if err != nil {
		logger.Panicf("error running server: %v", err)
	}
	defer func() {
		close(s.ServerError)
		close(s.ShutDown)
	}()
}
