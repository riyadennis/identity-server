package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/viper"

	"github.com/riyadennis/identity-server/internal"
	"github.com/riyadennis/identity-server/internal/handlers"
)

var (
	configFile = flag.String("config", "etc/config_test.yaml",
		"path to config file")
)

func init() {
	flag.Parse()
	// if we are not running main through docker-compose we
	// might not have environment variable set
	// then we take default from flag.
	cf := os.Getenv("CONFIG_FILE")
	if cf == "" {
		cf = *configFile
	}
	viper.SetConfigFile(cf)
	viper.AutomaticEnv()
	file, err := os.Create("identity.log")
	if err != nil {
		logrus.Fatalf("failed to create log file:: %v", err)
	}
	logrus.SetOutput(file)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("failed to read config :: %v", err)
	}
}

func main() {
	handlers.Init(viper.GetString("ENV"))
	internal.Server(viper.GetString("port"))
}
