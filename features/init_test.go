package features

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/cucumber/godog"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/internal/store"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
)

const HOST = "http://localhost:8095"

var (
	Idb    store.Store
)

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(beforeScenario)
	s.Step(`^a not registered email "([^"]*)"$`, aNotRegisteredEmail)
	s.Step(`^password "([^"]*)"$`, aNotRegisteredPassword)
	s.Step(`^I login$`, iLogin)
	s.Step(`^I should get error-code "([^"]*)"$`, iShouldGetErrorCode)
	s.Step(`status code (\d+)$`, statusCode)
	s.Step(`^message "([^"]*)"$`, message)

	s.Step(`^a registered user with email "([^"]*)"$`, aRegisteredUserWithEmail)
	s.Step(`^password "([^"]*)" with firstName "([^"]*)" and lastName "([^"]*)""$`, passwordWithFirstNameAndLastName)

	s.Step(`^that user login$`, thatUserLogin)
	s.Step(`^status code should be (\d+)$`, statusCode)

	s.AfterScenario(afterScenario)
}

func beforeScenario(f interface{}) {
	var err error
	Idb, err = connectSQLite()
	if err != nil {
		logrus.Fatal(err)
	}
	client = &http.Client{}
}

func afterScenario(i interface{}, e error) {
	err := sqlite.Truncate(viper.GetString("source"))
	if err != nil {
		logrus.Fatal(err)
	}
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
