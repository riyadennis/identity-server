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

	"github.com/riyadennis/identity-server/internal/store"
	"github.com/riyadennis/identity-server/internal/store/sqlM"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	Idb store.Store
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

func Init(env string) {
	if env == "test" {
		var err error
		Idb, err = connectSQLite()
		if err != nil {
			logrus.Fatal(err)
		}
		return
	}
	connectMysql()
}

func connectMysql() {
	var err error
	db, err := sqlM.ConnectDB()
	if err != nil {
		logrus.Fatal(err)
	}
	Idb, err = sqlM.PrepareDB(db)
	if err != nil {
		logrus.Fatal(err)
	}
}

func connectSQLite() (*sqlite.LiteDB, error) {
	db, err := sqlite.ConnectDB("/var/tmp/identity.db")
	if err != nil {
		return nil, err
	}
	err = sqlite.Setup("/var/tmp/identity.db")
	if err != nil {
		return nil, err
	}
	return sqlite.PrepareDB(db)
}

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
