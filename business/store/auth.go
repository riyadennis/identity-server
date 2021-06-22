package store

import "golang.org/x/crypto/bcrypt"

type Authenticator interface {
	Authenticate(email, password string) (bool, error)
}

// Authenticate checks the validity of a given password for an email
func (d *DB) Authenticate(email, password string) (bool, error) {
	login, err := d.Conn.Prepare(
		`SELECT  password FROM 
            	identity_users 
				where email = ?`)
	if err != nil {
		return false, err
	}

	row := login.QueryRow(email)
	var hashedPass string

	err = row.Scan(&hashedPass)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
