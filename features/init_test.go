package features

import (
	"net/http"

	"github.com/spf13/viper"

	_ "github.com/mattn/go-sqlite3"

	"github.com/riyadennis/identity-server/internal/store"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
)

const HOST = "http://localhost:8095"

var (
	client *http.Client
)

func beforeScenario(f interface{}) {
	client = &http.Client{}
}

func afterScenario(i interface{}, e error) {
}

func connectSQLite() (store.Store, error) {
	viper.SetConfigFile("../etc/config_test.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	db, err := sqlite.ConnectDB(viper.GetString("source"))
	if err != nil {
		return nil, err
	}
	err = sqlite.Setup(viper.GetString("source"))
	if err != nil {
		return nil, err
	}
	return store.PrepareDB(db)
}
