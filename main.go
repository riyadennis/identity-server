package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	mmysql "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/riyadennis/identity-server/internal"
)

var (
	port = flag.String("port", ":8080", "port http server will listen to")
)

const (
	//when we add a new migration this constant need to be updated
	step = 1
	//if we change the folder in which we keep our migration
	//files we need to update this
	sourceUrl = "file://migrations/"
)

func main() {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@/%s", user, password, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("could not connect to the MySQL database... %v", err)
	}
	driver, _ := mmysql.WithInstance(db, &mmysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		sourceUrl,
		"identity-server",
		driver,
	)

	m.Steps(step)
	flag.Parse()
	internal.Server(*port)
}
