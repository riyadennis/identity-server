package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Authenticate(email, password string) (bool, error)
	FetchLoginToken(userID string) (*TokenRecord, error)
	SaveLoginToken(ctx context.Context, t *TokenRecord) error
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

type TokenRecord struct {
	Id        string
	UserId    string
	Token     string
	TTL       string
	Expiry    time.Time
	LastUsed  sql.NullString
	CreatedAt string
	UpdatedAt string
}

var tokenQuery = `SELECT id,token,ttl,expiry,last_used FROM 
login_tokens 
where user_id = ?`

func (a *Auth) FetchLoginToken(userID string) (*TokenRecord, error) {
	query, err := a.Conn.Prepare(tokenQuery)
	if err != nil {
		return nil, err
	}
	token := &TokenRecord{}
	tokenRow := query.QueryRow(userID)
	err = tokenRow.Scan(&token.Id, &token.Token, &token.TTL, &token.Expiry, &token.LastUsed)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if token.Token == "" {
		a.Logger.Infof("no token found in DB new login")
		return nil, nil
	}
	return token, nil
}

var saveTokenQuery = `INSERT INTO login_tokens (id, user_id, token, ttl, expiry) VALUES (?, ?, ?, ?, ?)`

func (a *Auth) SaveLoginToken(ctx context.Context, t *TokenRecord) error {
	saveStmt, err := a.Conn.Prepare(saveTokenQuery)
	if err != nil {
		a.Logger.Errorf("failed to prepare save token query: %v", err)
		return err
	}
	id := uuid.New().String()
	result, err := saveStmt.ExecContext(ctx, id, t.UserId, t.Token, t.TTL, t.Expiry)
	if err != nil {
		a.Logger.Errorf("failed to save token: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		a.Logger.Errorf("failed to get rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return errors.New("failed to save token")
	}

	return nil
}
