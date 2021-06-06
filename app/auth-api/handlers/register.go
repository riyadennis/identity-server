package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

// Handler have common setup needed to run the handlers
// its helps to reuse open db connection
type Handler struct {
	Store  *store.DB
	Logger *log.Logger
}

func NewHandler(store *store.DB, logger *log.Logger) *Handler {
	return &Handler{
		Store:  store,
		Logger: logger,
	}
}

// Register is the handler function that will process
// rest call to register endpoint
func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := &store.UserRequest{}

	err := foundation.RequestBody(r, u)
	if err != nil {
		h.Logger.Printf("invalid data in request: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.InvalidRequest)
		return
	}

	err = validation.ValidateUser(u)
	if err != nil {
		h.Logger.Printf("validation failed:%v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.ValidationFailed)
		return
	}

	ctx := r.Context()
	exists, err := userExists(ctx, h.Store, u.Email)
	if err != nil {
		h.Logger.Println("failed to check user in database: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
		return
	}
	if exists {
		h.Logger.Println("user already exists: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, errors.New("email already exists"), foundation.EmailAlreadyExists)
		return
	}

	password, err := business.GeneratePassword()
	if err != nil {
		h.Logger.Println("failed to generate password: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.PassWordError)
		return
	}

	u.Password, err = business.EncryptPassword(password)
	if err != nil {
		h.Logger.Println("password encryption failed: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.PassWordError)
		return
	}

	resource, err := h.Store.Insert(ctx, u)
	if err != nil {
		h.Logger.Println("failed to save user: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
	}

	resource.Password = password

	_ = foundation.Resource(w, http.StatusCreated, resource)
}

func userDataFromRequest(r *http.Request) (*store.UserRequest, error) {
	u := &store.UserRequest{}
	err := foundation.RequestBody(r, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func userExists(ctx context.Context, store store.Store, email string) (bool, error) {
	selectUser, err := store.Read(ctx, email)
	if err != nil {
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
