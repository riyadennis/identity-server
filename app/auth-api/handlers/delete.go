package handlers

import (
	"errors"
	"github.com/riyadennis/identity-server/business/store"
	"net/http"

	"github.com/riyadennis/identity-server/foundation"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// UserDelete is the request structure for delete end point
// we need only email now.
type UserDelete struct {
	Email string `json:"email"`
}

// Delete is the handler for delete end point to remove a user from as per the email
func Delete(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	ud := &UserDelete{}
	err := foundation.RequestBody(req, ud)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			err, foundation.InvalidRequest)
		return
	}
	db, err := store.Connect()
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
	}

	done, err := store.NewDB(db).Delete(ud.Email)
	if err != nil {
		logrus.Errorf("failed to delete :: %v", err)
		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.InvalidRequest)
		return
	}
	if done == 0 {
		logrus.Error("failed to delete")
		foundation.ErrorResponse(w, http.StatusBadRequest, errors.New("failed to delete"), foundation.InvalidRequest)
		return
	}
	logrus.Infof("user %s deleted for :: %d records", ud.Email, done)
	err = foundation.JSONResponse(w, http.StatusOK,
		"account deleted",
		"")
	if err != nil {
		logrus.Error(err)
	}
}