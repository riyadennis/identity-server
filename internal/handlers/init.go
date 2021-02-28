package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/riyadennis/identity-server/internal/store"
)

var (
	Idb store.Store
)

const (
	passwordSeed = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tokenTTL     = 120 * time.Hour
)

type Response struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"error-code"`
}

// NewCustomError returns error with error code
func NewCustomError(code string, err error) *CustomError {
	return &CustomError{
		Code: code,
		Err:  err,
	}
}

func dataSource() store.Store {
	return Idb
}

func encryptPassword(password string) (string, error) {
	enPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(enPass), nil

}

func generatePassword() (string, error) {
	result := ""
	for {
		if len(result) >= 15 {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		s := fmt.Sprintf("%d", num.Int64())
		if strings.Contains(passwordSeed, s) {
			result += s
		}
	}
}

func requestBody(r *http.Request) ([]byte, error) {
	if r.Header.Get("content-type") != "application/json" {
		err := errors.New("invalid content type")
		logrus.Errorf("content type is not json :: %v", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if err == io.EOF {
			err := errors.New("empty request body")
			return nil, err
		}
		return nil, err
	}
	return data, nil
}
