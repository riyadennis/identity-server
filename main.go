package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
)

func init() {
	if os.Getenv("ENV") != "test" {
		file, err := os.Create("identity.log")
		if err != nil {
			logrus.Fatalf("failed to create log file:: %v", err)
		}
		logrus.SetOutput(file)
	}
}

func main() {
	if os.Getenv("ENV") == "" || os.Getenv("PORT") == "" {
		logrus.Fatal("invalid setting no port number")
	}
	handlers.Init(os.Getenv("ENV"))
	internal.Server(os.Getenv("PORT"))
}
