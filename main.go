package main

import (
	"flag"

	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
)

var (
	port = flag.String("port", ":8095", "port http server will listen to")
)

func main() {
	flag.Parse()
	handlers.Init()
	internal.Server(*port)
}
