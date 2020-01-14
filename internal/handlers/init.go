package handlers

import (
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/riyadennis/identity-server/internal/store/sqlite"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	Idb *sqlite.LiteDB
)

const (
	passwordSeed = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tokenTTL     = 120 * time.Second
	mySigningKey = "thisistobereplaced"
)

type Response struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"error-code"`
	Token     string `json:"token,omitempty"`
}

func Init() {
	Idb = sqlite.PrepareDB("/var/tmp/identity.db")
}

func NewCustomError(code string, err error) *CustomError {
	return &CustomError{
		Code: code,
		Err:  err,
	}
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
		s := string(num.Int64())
		if strings.Contains(passwordSeed, s) {
			result += s
		}
	}
}

func requestBody(r *http.Request) ([]byte, error) {
	if r.Header.Get("content-type") != "application/json" {
		err := errors.New("invalid content")
		logrus.Error(err)
		return nil, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if err == io.EOF {
			err := errors.New("invalid content")
			logrus.Error(err)
			return nil, err
		}
		return nil, err
	}
	return data, nil
}
