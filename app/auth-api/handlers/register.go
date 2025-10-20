package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/julienschmidt/httprouter"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

// Handler have common setup needed to run the handlers
// its helps to reuse open db connection
type Handler struct {
	Store         store.Store
	Authenticator store.Authenticator
	Logger        *log.Logger
	TokenConfig   *store.TokenConfig
}

func NewHandler(store store.Store, authenticator store.Authenticator,
	tc *store.TokenConfig, logger *log.Logger) *Handler {
	return &Handler{
		Store:         store,
		Authenticator: authenticator,
		Logger:        logger,
		TokenConfig:   tc,
	}
}

// @Summary      Register a new user
// @Description  Create a user with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body   store.UserRequest  true  "User registration data"
// @Success      201   {object}  store.UserResource
// @Failure      400   {object}  foundation.Response
// @Failure      500   {object}  foundation.Response
// @Router       /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := &store.UserRequest{}

	err := jsonapi.UnmarshalPayload(r.Body, u)
	if err != nil {
		h.Logger.Printf("invalid data in request: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.InvalidRequest)
		return
	}

	err = validation.ValidateUser(u)
	if err != nil {
		h.Logger.Printf("validation failed: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.ValidationFailed)
		return
	}

	userExists, err := h.Store.Read(r.Context(), u.Email)
	if err != nil {
		h.Logger.Printf("failed to read from database: %v", err)

		foundation.ErrorResponse(w, http.StatusBadRequest, err, foundation.ValidationFailed)
		return
	}

	if userExists.Email == u.Email {
		h.Logger.Printf("email already exists: %#v", userExists.Email)

		foundation.ErrorResponse(w, http.StatusBadRequest, errors.New("email already exists"), foundation.EmailAlreadyExists)
		return
	}

	u.Password, err = business.EncryptPassword(u.Password)
	if err != nil {
		h.Logger.Printf("password encryption failed: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.PassWordError)
		return
	}

	resource, err := h.Store.Insert(r.Context(), u)
	if err != nil {
		h.Logger.Printf("failed to save user: %v", err)

		foundation.ErrorResponse(w, http.StatusInternalServerError, err, foundation.DatabaseError)
	}

	resource.Password = "********"
	_ = foundation.Resource(w, http.StatusCreated, resource)
}
