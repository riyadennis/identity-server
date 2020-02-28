package features

import (
	"net/http"

	"github.com/cucumber/godog"
	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
	"github.com/sirupsen/logrus"
)

const HOST = "http://localhost:8088"

var (
	Idb    store.Store
	client *http.Client
)

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(beforeScenario)

	s.Step(`^email "([^"]*)"$`, email)
	s.Step(`^password "([^"]*)"$`, password)
	s.Step(`^I login$`, iLogin)
	s.Step(`^I should get error-code "([^"]*)"$`, iShouldGetErrorcode)
	s.Step(`status code (\d+)$`, statusCode)
	s.Step(`^message "([^"]*)"$`, message)

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
	//TODO to truncate db
}

func connectSQLite() (*sqlite.LiteDB, error) {
	db, err := sqlite.ConnectDB("/var/tmp/identity.db")
	if err != nil {
		return nil, err
	}
	err = sqlite.Setup("/var/tmp/identity.db")
	if err != nil {
		return nil, err
	}
	return sqlite.PrepareDB(db)
}
