package store

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Authenticate(email, password string) (bool, error)
}

var authQuery = `SELECT password FROM 
identity_users 
where email = ?`

// Authenticate checks the validity of a given password for an email
func (d *DB) Authenticate(email, inputPassword string) (bool, error) {
	login, err := d.Conn.Prepare(authQuery)
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
		logrus.Errorf("hashed password error :: %v", err)
		return false, err
	}

	return true, nil
}
