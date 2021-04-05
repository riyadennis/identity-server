package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

// Handler have common setup needed to run the handlers
// its helps to reuse open db connection
type Handler struct {
	Store store.Store
}

// Register is the handler function that will process
// rest call to register endpoint
func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, err := userDataFromRequest(r)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.InvalidRequest)
		return
	}

	err = validation.ValidateUser(u)
	if err != nil {
		logrus.Errorf("validation failed :: %v", err)
		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.ValidationFailed)
		return
	}

	if h.Store == nil {
		logrus.Error("invalid db connection")
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errors.New("invalid db connection"), foundation.DatabaseError)
		return
	}

	ctx := r.Context()
	exists, err := userExists(ctx, h.Store, u.Email)
	if err != nil {
		// already logged
		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
		return
	}
	if exists {
		foundation.ErrorResponse(w, http.StatusBadRequest, errors.New("email already exists"), foundation.EmailAlreadyExists)
		return
	}

	password, err := business.GeneratePassword()
	if err != nil {
		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.PassWordError)
		return
	}

	u.Password, err = business.EncryptPassword(password)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.PassWordError)
		return
	}

	err = storeUser(ctx, h.Store, u)
	if err != nil {
		logrus.Errorf("failed to save user  :: %v", err)
		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = foundation.JSONResponse(w, http.StatusOK,
		fmt.Sprintf("your generated password : %s", password),
		"")
	if err != nil {
		logrus.Error(err)
	}
}

func userDataFromRequest(r *http.Request) (*store.User, error) {
	if r == nil {
		return nil, errors.New("empty request")
	}
	reqID := r.Context().Value("reqID")
	u := &store.User{}
	err := foundation.RequestBody(r, u)
	if err != nil {
		logrus.Errorf("requestID %s failed to read request body :: %v", reqID, err)
		return nil, err
	}

	return u, nil
}

func storeUser(ctx context.Context, store store.Store, u *store.User) error {
	err := store.Insert(ctx, u, uuid.New().String())
	if err != nil {
		logrus.Errorf("failed to register :: %v", err)
		return err
	}
	return nil
}

func userExists(ctx context.Context, store store.Store, email string) (bool, error) {
	selectUser, err := store.Read(ctx, email)
	if err != nil {
		logrus.Errorf("failed to check for user in database :: %v", err)
		return false, err
	}

	if selectUser != nil {
		if selectUser.Email == email {
			return true, nil
		}

		return false, nil
	}
	return false, nil
}
