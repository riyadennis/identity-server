package store

import (
	"database/sql"
	"fmt"
	"os"
)

type Config struct {
	BasePath string
	DB       DBConnection
}
type DBConnection struct {
	User      string
	Password  string
	Host      string
	Name      string
	Database  string
	Port      string
	ParseTime bool
}

func NewENVConfig() Config {
	return Config{
		BasePath: os.Getenv("BASE_PATH"),
		DB: DBConnection{
			User:      os.Getenv("MYSQL_USERNAME"),
			Password:  os.Getenv("MYSQL_PASSWORD"),
			Host:      os.Getenv("MYSQL_HOST"),
			Database:  os.Getenv("MYSQL_DATABASE"),
			Port:      os.Getenv("MYSQL_PORT"),
			ParseTime: true,
		},
	}
}

// Connect opens a connection to mysql
func Connect(cfg DBConnection) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.ParseTime,
	)

	return sql.Open("mysql", dsn)
}
