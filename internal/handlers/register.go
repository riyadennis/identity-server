package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
)

// Register is the handler function that will process rest call to register endpoint
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := requestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	u := &store.User{}
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
	exists, err := userExists(u.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("failed to check user  :: %v", err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	if exists {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "email already exists :: %v", u.Email)
		return
	}
	password, err := generatePassword()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to generate password :: %v", u.Email)
		return
	}
	u.Password = password
	err = storeUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("failed to save user  :: %v", err)
		fmt.Fprintf(w, "%v", err)
	}
	w.Write([]byte(password))
}

func storeUser(u *store.User) error {
	err := Idb.Insert(u)
	if err != nil {
		logrus.Errorf("failed to register :: %v", err)
		return err
	}
	return nil
}

func userExists(email string) (bool, error) {
	selectUser, err := Idb.Read(email)
	if err != nil {
		logrus.Errorf("failed to check for user in database :: %v", err)
		return false, err
	}
	if selectUser == nil {
		return false, nil
	}
	return true, nil
}

func validateUser(u *store.User) error {
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
