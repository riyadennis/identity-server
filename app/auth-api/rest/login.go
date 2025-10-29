package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/validation"
	"github.com/riyadennis/identity-server/foundation"
)

var (
	errEmailNotFound   = errors.New("email not found")
	errTokenGeneration = errors.New("key not found")
)

// UserLogin have data needed for a user to login
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login @Summary      Login Endpoint
//
//	@Description	Authenticate a user and return a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Basic base64(email:password)"
//	@Success		200				{object}	store.Token
//	@Failure		400				{object}	foundation.Response
//	@Failure		401				{object}	foundation.Response
//	@Failure		500				{object}	foundation.Response
//	@Router			/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email, password, ok := r.BasicAuth()
	if !ok {
		h.Logger.Printf("invalid request, username or password is empty")

		foundation.ErrorResponse(w, http.StatusBadRequest,
			errors.New("empty login data"), foundation.InvalidRequest)
		return
	}
	err := validation.ValidateEmail(email)
	if err != nil {
		h.Logger.Printf("invalid request, username is invalid: %v", err)
		foundation.ErrorResponse(w, http.StatusBadRequest,
			err, foundation.InvalidRequest)
		return
	}
	helper := business.NewHelper(h.Store, h.Authenticator, h.Logger)
	user, err := helper.UserCredentialsInDB(r.Context(), email, password)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusBadRequest,
			err, foundation.InvalidRequest)
		return
	}

	token, err := helper.ManageToken(r.Context(), h.TokenConfig, user.ID)
	if err != nil {
		foundation.ErrorResponse(w, http.StatusInternalServerError,
			errTokenGeneration, foundation.TokenError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		h.Logger.Printf("json encoding failed: %v", err)
	}
}
