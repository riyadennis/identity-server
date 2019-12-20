package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/riyadennis/identity-server/internal/store"
	"github.com/sirupsen/logrus"
)

// Register is the handler function that will process rest call to register endpoint
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	data, err := requestBody(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}

	u := &store.User{}
	err = json.Unmarshal(data, u)
	if err != nil {
		logrus.Errorf("failed to unmarshal :: %v", err)
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: InvalidRequest,
			Err:  err,
		})
		return
	}
	customErr := validateUser(u)
	if customErr != nil {
		logrus.Errorf("validation failed :: %v", customErr.Err)
		errorResponse(w, http.StatusBadRequest, customErr)
		return
	}

	password, err := generatePassword()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, &CustomError{
			Code: PassWordError,
			Err:  err,
		})
		return
	}
	enPass, err := encryptPassword(password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: PassWordError,
			Err:  err,
		})
		return
	}

	u.Password = enPass
	err = storeUser(u)
	if err != nil {
		logrus.Errorf("failed to save user  :: %v", err)
		errorResponse(w, http.StatusInternalServerError, &CustomError{
			Code: DatabaseError,
			Err:  err,
		})
	}

	err = jsonResponse(w, http.StatusOK,
		fmt.Sprintf("your generated password : %s", password),
		"")
	if err != nil {
		logrus.Error(err)
	}
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

func validateUser(u *store.User) *CustomError {
	if u.FirstName == "" {
		return &CustomError{
			Code: FirstNameMissing,
			Err:  errors.New("missing first name"),
		}
	}
	if u.LastName == "" {
		return &CustomError{
			Code: LastNameMissing,
			Err:  errors.New("missing last name"),
		}
	}
	if u.Email == "" {
		return &CustomError{
			Code: EmailMissing,
			Err:  errors.New("missing last name"),
		}
	}
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(u.Email) {
		return &CustomError{
			Code: EmailInvalid,
			Err:  errors.New("invalid email"),
		}
	}
	if u.Terms == false {
		return &CustomError{
			Code: TermsMissing,
			Err:  errors.New("missing terms"),
		}
	}
	exists, err := userExists(u.Email)
	if err != nil {
		return &CustomError{
			Code: DatabaseError,
			Err:  err,
		}
	}
	if exists {
		return &CustomError{
			Code: EmailExists,
			Err:  fmt.Errorf("email already exists :: %v", u.Email),
		}
	}
	return nil
}
