package main

import (
	"flag"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
)

var (
	configFile = flag.String("config", "etc/config_test.yaml",
		"path to config file")
)

func main() {
	flag.Parse()
	viper.SetConfigFile(*configFile)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("failed to read config :: %v", err)
	}
	handlers.Init(viper.GetString("ENV"))
	internal.Server(viper.GetString("port"))
}
