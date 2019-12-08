package main

import (
	"flag"

	"github.com/riyadennis/identity-server/internal"
)

var (
	port = flag.String("port", ":8087", "port http server will listen to")
)

func main() {
	flag.Parse()
	internal.Server(*port)
}
