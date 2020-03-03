package features

import (
	"github.com/cucumber/godog"
	_ "github.com/mattn/go-sqlite3"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
	"github.com/sirupsen/logrus"
	"net/http"
)

const HOST = "http://localhost:8088"

var (
	Idb    *sqlite.LiteDB
	client *http.Client
)

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(beforeScenario)
	s.Step(`^a not registered email "([^"]*)"$`, aNotRegisteredEmail)
	s.Step(`^password "([^"]*)"$`, password)
	s.Step(`^I login$`, iLogin)
	s.Step(`^I should get error-code "([^"]*)"$`, iShouldGetErrorcode)
	s.Step(`status code (\d+)$`, statusCode)
	s.Step(`^message "([^"]*)"$`, message)

	s.Step(`^a registered user with email "([^"]*)" with firstName "([^"]*)" and lastName "([^"]*)""$`, aRegisteredUserWithEmailWithFirstNameAndLastName)
	s.Step(`^that user login$`, thatUserLogin)
	s.Step(`^status code should be (\d+)$`, statusCode)
	s.Step(`^token not "([^"]*)"$`, tokenNot)


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
