package main

import (
	"os"
	"os/signal"
	"syscall"

	// initialise mysql driver
	// initialise migration settings
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/riyadennis/identity-server/app/gql/graph"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func main() {
	logger := foundation.NewLogger()
	err := foundation.ValidatePort(os.Getenv("GRAPHQL_PORT"))
	if err != nil {
		logger.Fatalf("invalid port: %v", os.Getenv("GRAPHQL_PORT"))
	}
	cfg := store.NewENVConfig()
	st, auth, err := store.SetUpMYSQL(logger)
	if err != nil {
		logger.Fatalf("database setUp failed %v", err)
	}
	s := graph.NewServer(logger, os.Getenv("GRAPHQL_PORT"), st, auth, cfg.Token)
	signal.Notify(s.ShutDown, os.Interrupt, syscall.SIGTERM)

	err = s.Start(os.Getenv("GRAPHQL_PORT"))
	if err != nil {
		logger.Fatalf("failed to start graphQL server %v", err)
	}

	<-s.ShutDown

}
