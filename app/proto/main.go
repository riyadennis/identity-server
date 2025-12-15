package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	// initialise mysql driver
	// initialise migration settings
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/riyadennis/identity-server/app/proto/identity"

	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/foundation"
)

func main() {
	logger := foundation.NewLogger()
	err := foundation.ValidatePort(os.Getenv("GRPC_PORT"))
	if err != nil {
		logger.Fatalf("invalid port: %v", os.Getenv("GRPC_PORT"))
	}
	cfg := store.NewENVConfig()
	st, auth, err := store.SetUpMYSQL(logger)
	if err != nil {
		logger.Fatalf("database setUp failed %v", err)
	}
	server := identity.NewServer(logger, cfg.Token, st, auth)
	signal.Notify(server.ShutDown, os.Interrupt, syscall.SIGTERM)
	var wt sync.WaitGroup
	go func() {
		wt.Add(1)
		err = server.Run(os.Getenv("GRPC_PORT"))
		if err != nil {
			logger.Fatalf("failed to run gRPC server: %v", err)
			wt.Done()
		}
	}()
	<-server.ShutDown
	wt.Wait()
}
