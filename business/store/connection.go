package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var (
	errInvalidDBConfig = errors.New("invalid db configuration")
	errEmptyDBUserName = errors.New("empty mysql user name")
	errEmptyDBPassword = errors.New("empty mysql password")
	errEmptyDBHost     = errors.New("empty mysql host name")
	errEmptyDBName     = errors.New("empty mysql database name")
	errEmptyDBPort     = errors.New("empty mysql port")
)

type Config struct {
	DB    *DBConnection
	Token *TokenConfig
}

type TokenConfig struct {
	Issuer         string
	KeyPath        string
	PrivateKeyName string
	PublicKeyName  string
	TokenTTL       time.Duration
}

type DBConnection struct {
	User          string
	Password      string
	Host          string
	Name          string
	Database      string
	Port          string
	ParseTime     bool
	MigrationPath string
}

// Token have credentials present in a token
type Token struct {
	Status      int           `json:"status"`
	AccessToken string        `json:"access_token"`
	Expiry      string        `json:"expiry"`
	TokenType   string        `json:"token_type"`
	TokenTTL    time.Duration `json:"token_ttl" swaggertype:"string"`
}

func NewENVConfig() *Config {
	return &Config{
		DB: &DBConnection{
			User:          os.Getenv("MYSQL_USERNAME"),
			Password:      os.Getenv("MYSQL_PASSWORD"),
			Host:          os.Getenv("MYSQL_HOST"),
			Database:      os.Getenv("MYSQL_DATABASE"),
			Port:          os.Getenv("MYSQL_PORT"),
			ParseTime:     true,
			MigrationPath: os.Getenv("MIGRATION_PATH"),
		},
		Token: &TokenConfig{
			Issuer:         os.Getenv("ISSUER"),
			KeyPath:        os.Getenv("KEY_PATH"),
			PrivateKeyName: "private.pem",
			PublicKeyName:  "public.pem",
		},
	}
}

// Connect opens a connection to mysql
func Connect(dbCfg *DBConnection) (*sql.DB, error) {
	if dbCfg == nil {
		return nil, errInvalidDBConfig
	}
	if dbCfg.User == "" {
		return nil, errEmptyDBUserName
	}
	if dbCfg.Password == "" {
		return nil, errEmptyDBPassword
	}

	if dbCfg.Host == "" {
		return nil, errEmptyDBHost
	}

	if dbCfg.Port == "" {
		return nil, errEmptyDBPort
	}

	if dbCfg.Database == "" {
		return nil, errEmptyDBName
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database,
		dbCfg.ParseTime,
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GenerateToken(logger *logrus.Logger, issuer string, key []byte, ttl time.Duration) (*Token, error) {
	expiry := time.Now().UTC().Add(ttl)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		logger.Printf("failed to parser private key: %v", err)
		return nil, err
	}

	t, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": expiry.Unix(),
		"iss": issuer,
	}).SignedString(privateKey)
	if err != nil {
		logger.Errorf("failed to sign using private key: %v", err)
		return nil, err
	}

	return &Token{
		Status:      200,
		AccessToken: t,
		Expiry:      expiry.String(),
		TokenType:   "Bearer",
		TokenTTL:    ttl,
	}, nil
}
