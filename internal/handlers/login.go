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
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}
	var ld *LoginDetails
	if data != nil {
		ld = &LoginDetails{}
		err = json.Unmarshal(data, ld)
		if err != nil {
			logrus.Errorf("failed to unmarshal :: %v", err)
			errorResponse(w, http.StatusBadRequest, &CustomError{
				Code: InvalidRequest,
				Err:  err,
			})
			return
		}
	}

	if ld.Email == "" {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: EmailMissing,
			Err:  errors.New("email missing"),
		})
		return
	}
	if ld.Password == "" {
		logrus.Error("no password")
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: PassWordError,
			Err:  errors.New("password missing"),
		})
		return
	}
	fName, err := Idb.Authenticate(ld.Email, ld.Password)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResponse(http.StatusOK,
		fmt.Sprintf("welcome  : %s", fName),
		"",
	))
	return
}
