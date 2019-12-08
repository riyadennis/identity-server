package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// User hold information needed to complete user registration
type User struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Company          string `json:"company"`
	PostCode         string `json:"post_code"`
	Terms            bool   `json:"terms"`
	RegistrationDate string
}

// Register is the handler function that will process rest call to register endpoint
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid content type"))
		logrus.Error(errors.New("invalid content"))
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
		logrus.Errorf("failed to unmarshal :: %v", err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	err = validateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Errorf("validation failed :: %v", err)
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func validateUser(u *User) error {
	if u.FirstName == "" {
		return errors.New("missing first name")
	}
	if u.LastName == "" {
		return errors.New("missing last name")
	}
	if u.Email == "" {
		return errors.New("missing email")
	}
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(u.Email) {
		return errors.New("invalid email")
	}
	if u.Terms == false {
		return errors.New("missing terms")
	}
	return nil
}
