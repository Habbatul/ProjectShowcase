package config

import (
	"database/sql"
	"log"
)

var db *sql.DB

func ConnectDB() *sql.DB {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=testfirstgo_db password=mysecretpassword sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
