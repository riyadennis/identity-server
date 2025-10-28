package store

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Authenticate(email, password string) (bool, error)
}

type Auth struct {
	Conn   *sql.DB
	Logger *logrus.Logger
}

var authQuery = `SELECT password FROM 
identity_users 
where email = ?`

// Authenticate checks the validity of a given password for an email
func (a *Auth) Authenticate(email, inputPassword string) (bool, error) {
	login, err := a.Conn.Prepare(authQuery)
	if err != nil {
		return false, err
	}

	row := login.QueryRow(email)
	var storedHash string

	err = row.Scan(&storedHash)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword))
	if err != nil {
		a.Logger.Errorf("hashed password error :: %v", err)
		return false, err
	}

	return true, nil
}
