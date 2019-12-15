package handlers

import (
	"encoding/json"
	"errors"
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
		logrus.Errorf("failed to unmarshal :: %v", err)
		errorResponse(w, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}
	if ld.Email == "" {
		logrus.Errorf("no email :: %v", ld)
		errorResponse(w, &CustomError{
			Code: EmailMissing,
			Err:  err,
		})
		return
	}
	if ld.Password == "" {
		logrus.Error("no password")
		errorResponse(w, &CustomError{
			Code: PassWordError,
			Err:  err,
		})
		return
	}
	fname, err := Idb.Authenticate(ld.Email, ld.Password)
	if err != nil {
		errorResponse(w, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}
	if fname == "" {
		errorResponse(w, &CustomError{
			Code: InvalidRequest,
			Err:  errors.New("cannot authenticate"),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResponse(http.StatusOK,
		fmt.Sprintf("welcome  : %s", fname),
		"",
	))
	return
}
