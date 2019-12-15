package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	data, err := requestBody(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	ld := &LoginDetails{}
	err = json.Unmarshal(data, ld)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Errorf("failed to unmarshal :: %v", err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	if ld.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Errorf("no email :: %v", ld)
		fmt.Fprint(w, "no email")
		return
	}
	if ld.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Error("no password")
		fmt.Fprint(w, "no password")
		return
	}
	a, err := Idb.Authenticate(ld.Email, ld.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "unable to authenticate the user")
		return
	}
	if !a {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "unable to login")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "welcome")
	return
}
