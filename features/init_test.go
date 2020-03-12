package features

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const HOST = "http://localhost:8095"

func beforeScenario(f interface{}) {
	client = &http.Client{}
}

func afterScenario(i interface{}, e error) {
}
