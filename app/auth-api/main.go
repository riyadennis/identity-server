package main

import (
	"os"
	"os/signal"
	"syscall"

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
		logger.Fatalf("migration failed: %v", err)
	}

	s, err := server.NewServer(logger, os.Getenv("PORT"))
	if err != nil {
		logger.Fatalf("server initialisation failed: %v", err)
	}
	signal.Notify(s.ShutDown, os.Interrupt, syscall.SIGTERM)

	defer func() {
		close(s.ServerError)
		close(s.ShutDown)
	}()

	s.RESTHandler(db, cfg.Token)
	err = s.Run()
	if err != nil {
		logger.Fatalf("error running server: %v", err)
	}

}
