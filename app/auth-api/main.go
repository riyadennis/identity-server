package main

import (
	"os"
	"os/signal"
	"syscall"

	// initialise mysql driver
	// initialise migration settings
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/riyadennis/identity-server/app/auth-api/server"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func main() {
	logger := foundation.NewLogger()
	cfg := store.NewENVConfig()

	st, auth, err := store.SetUpMYSQL(logger)
	if err != nil {
		logger.Fatalf("database setUp failed %v", err)
	}

	newServer, err := server.NewServer(logger, os.Getenv("REST_PORT"))
	if err != nil {
		logger.Fatalf("server initialisation failed: %v", err)
	}
	signal.Notify(newServer.ShutDown, os.Interrupt, syscall.SIGTERM)

	defer func() {
		close(newServer.ServerError)
		close(newServer.ShutDown)
	}()

	newServer.RESTHandler(cfg.Token, st, auth)

	err = newServer.Run()
	if err != nil {
		logger.Fatalf("error running server: %v", err)
	}

}
