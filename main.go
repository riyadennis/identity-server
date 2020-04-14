package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/viper"

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
	handlers.Init(viper.GetString("ENV"))
	internal.Server(viper.GetString("PORT"))
}
