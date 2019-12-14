package main

import (
	"flag"

	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
)

var (
	port = flag.String("port", ":8080", "port http server will listen to")
)

func main() {
	flag.Parse()
	err := sqlite.Setup()
	if err != nil {
		panic(err)
	}
	internal.Server(*port)
}
