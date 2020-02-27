package main

import (
	"flag"
	"os"

	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
)

var (
	port = flag.String("port", ":8095",
		"port http server will listen to")
)

func main() {
	flag.Parse()
	handlers.Init(os.Getenv("ENV"))
	internal.Server(*port)
}
