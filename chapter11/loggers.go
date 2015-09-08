package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"log"
	"os"
)

var (
	Database       *log.Logger
	Authentication *log.Logger
	Errors         *log.Logger
)

func LogPrepare() {
	dblog, err := os.OpenFile("database.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	authlog, err := os.OpenFile("auth.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	errlog, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}

	Database = log.New(dblog, "DB:", log.Ldate|log.Ltime)
	Authentication = log.New(authlog, "AUTH:", log.Ldate|log.Ltime)
	Errors = log.New(errlog, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	LogPrepare()

	db, err := sql.Open("mysql", "root@/localhost")
	if err != nil {

	}

	d, _ := db.Prepare("SELECT * FROM content_stories")

	Database.Println(d)
	Authentication.Println("Logging an auth attempt item")
	Errors.Println("Logging an error")

}
