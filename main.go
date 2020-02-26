package main

import (
	"flag"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
)

var (
	port = flag.String("port", ":8081", "port http server will listen to")
)

func main() {
	flag.Parse()
	handlers.Init()
	internal.Server(*port)
}
