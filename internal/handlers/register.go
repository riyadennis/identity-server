package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// User hold information needed to complete user registration
type User struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Company          string `json:"company"`
	PostCode         string `json:"post_cod"`
	Terms            bool   `json:"terms"`
	RegistrationDate string
}

// Register is the handler function that will process rest call to register endpoint
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid content type"))
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == io.EOF {
			w.Write([]byte("empty request body"))
			return
		}
		w.Write([]byte(err.Error()))
		return
	}
	u := &User{}
	err = json.Unmarshal(data, u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

