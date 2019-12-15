package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
)

var (
	port = flag.String("port", ":8080", "port http server will listen to")
)

func main() {
	flag.Parse()
	err := sqlite.Setup("/var/tmp/identity.db")
	if err != nil {
		logrus.Errorf("failed to load db :: %v", err)
		os.Exit(2)
	}
	handlers.Init()
	internal.Server(*port)
}
