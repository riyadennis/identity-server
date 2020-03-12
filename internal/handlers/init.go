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

	"github.com/spf13/viper"

	"github.com/riyadennis/identity-server/internal/store"
	"github.com/riyadennis/identity-server/internal/store/sqlM"
	"github.com/riyadennis/identity-server/internal/store/sqlite"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	Idb *store.DB
)

const (
	passwordSeed = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tokenTTL     = 120 * time.Second
)

type Response struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"error-code"`
}

// Init initialises and loads database settings
func Init(env string) {
	var err error
	// if environment is test
	// we want to initialise sqlite
	// database.

	if env == "test" {
		Idb, err = connectSQLite()
		if err != nil {
			logrus.Fatal(err)
		}
		return
	}
	Idb, err = connectMysql()
	if err != nil {
		logrus.Fatal(err)
	}
}

// NewCustomError returns error with error code
func NewCustomError(code string, err error) *CustomError {
	return &CustomError{
		Code: code,
		Err:  err,
	}
}

func connectMysql() (*store.DB, error) {
	var err error
	db, err := sqlM.ConnectDB()
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}
	logrus.Infof("MYSQL db details %v", db.Stats())
	return store.PrepareDB(db)
}

func connectSQLite() (*store.DB, error) {
	db, err := sqlite.ConnectDB(viper.GetString("source"))
	if err != nil {
		return nil, err
	}
	err = sqlite.Setup(viper.GetString("source"))
	if err != nil {
		return nil, err
	}
	logrus.Infof("SQLite db details %#v", db.Stats().MaxOpenConnections)
	return store.PrepareDB(db)
}

func dataSource() *store.DB {
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
