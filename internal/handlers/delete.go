package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type UserDelete struct {
	Email string `json:"email"`
}

func Delete(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	data, err := requestBody(req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidRequest, err))
		return
	}
	ud := &UserDelete{}
	err = json.Unmarshal(data, ud)
	if err != nil {
		logrus.Errorf("failed to unmarshal :: %v", err)
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidRequest, err))
		return
	}
	done, err := Idb.Delete(ud.Email)
	if err != nil {
		logrus.Errorf("failed to delete :: %v", err)
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidRequest, err))
		return
	}
	if done == 0 {
		logrus.Error("failed to delete")
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidRequest, err))
		return
	}
	logrus.Infof("user %s deleted for :: %d records", ud.Email, done)
	err = jsonResponse(w, http.StatusOK,
		"account deleted",
		"")
	if err != nil {
		logrus.Error(err)
	}

}
