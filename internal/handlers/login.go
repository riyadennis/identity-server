package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var mySigningKey = []byte("captainjacksparrowsayshi")

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
	token, err := generateToken()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: TokenError,
			Err:  err,
		})
		return
	}
	res  := newResponse(http.StatusOK,
		fmt.Sprintf("welcome  : %s", fName),
		"",
	)
	res.Token = token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return
}

func generateToken() (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil{
		return "", err
	}
	return tokenStr, nil
}