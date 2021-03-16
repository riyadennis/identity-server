package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
)

// Register is the handler function that will process
// rest call to register endpoint
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r == nil {
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidRequest, errors.New("empty request")))
		return
	}
	ctx := r.Context()
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u, err := userDataFromRequest(cctx, r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidRequest, err))
		return
	}
	err = validateUser(u)
	if err != nil {
		logrus.Errorf("validation failed :: %v", err)
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(InvalidUserData, err))
		return
	}
	exists, err := userExists(u.Email)
	if err != nil {
		//already logged
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(DatabaseError, err))
		return
	}
	if exists {
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(EmailAlreadyExists,
				errors.New("email already exists")))
		return
	}

	password, err := generatePassword()
	if err != nil {
		errorResponse(w, http.StatusBadRequest,
			NewCustomError(PassWordError, err))
		return
	}
	u.Password, err = encryptPassword(password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError,
			NewCustomError(PassWordError, err))
		return
	}

	err = storeUser(u)
	if err != nil {
		logrus.Errorf("failed to save user  :: %v", err)
		errorResponse(w, http.StatusInternalServerError,
			NewCustomError(DatabaseError, err))
	}
	w.Header().Set("Content-Type", "application/json")
	err = jsonResponse(w, http.StatusOK,
		fmt.Sprintf("your generated password : %s", password),
		"")
	if err != nil {
		logrus.Error(err)
	}
}

func userDataFromRequest(ctx context.Context, r *http.Request) (*store.User, error) {
	reqID := ctx.Value("reqID")
	if r == nil {
		return nil, errors.New("empty request")
	}
	data, err := requestBody(r)
	if err != nil {
		logrus.Errorf("requestID %s failed to read request body :: %v", reqID, err)
		return nil, err
	}
	u := &store.User{}
	err = json.Unmarshal(data, u)
	if err != nil {
		logrus.Errorf("failed to unmarshal :: %v", err)
		return nil, err
	}
	return u, nil
}

func errorResponse(w http.ResponseWriter, code int, customErr *CustomError) {
	w.Header().Set("Content-Type", "application/json")
	err := jsonResponse(w, code, customErr.Error(), customErr.Code)
	if err != nil {
		logrus.Error(err)
	}
}

func jsonResponse(w http.ResponseWriter, status int,
	message, errCode string) error {
	w.WriteHeader(status)
	res := newResponse(status, message, errCode)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func newResponse(status int, message, errCode string) *Response {
	return &Response{
		Status:    status,
		Message:   message,
		ErrorCode: errCode,
	}
}

func storeUser(u *store.User) error {
	err := Idb.Insert(u)
	if err != nil {
		logrus.Errorf("failed to register :: %v", err)
		return err
	}
	return nil
}

func userExists(email string) (bool, error) {
	selectUser, err := Idb.Read(email)
	if err != nil {
		logrus.Errorf("failed to check for user in database :: %v", err)
		return false, err
	}
	if selectUser == nil {
		return false, nil
	}
	return true, nil
}

func validateUser(u *store.User) error {
	if u == nil {
		return errors.New("empty user details")
	}
	if u.FirstName == "" {
		return errors.New("missing first name")
	}
	if u.LastName == "" {
		return errors.New("missing last name")
	}
	if u.Email == "" {
		return errors.New("missing email")
	}
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(u.Email) {
		return errors.New("invalid email")
	}
	if u.Terms == false {
		return errors.New("missing terms")
	}
	return nil
}
