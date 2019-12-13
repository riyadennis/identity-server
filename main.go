package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/riyadennis/identity-server/internal"
)

var (
	port = flag.String("port", ":8080", "port http server will listen to")
)


func main() {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@@tcp(172.17.0.2:3306)/%s", user, password, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, err = conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS identity_users(id varchar(100) NOT NULL PRIMARY KEY,first_name  varchar(100),last_name varchar(100),email varchar(100),company varchar(100),post_code varchar(100),terms int, created_datetime DATETIME)")
	if err != nil {
		log.Fatalf("%v", err)
	}
	flag.Parse()
	internal.Server(*port)
}
