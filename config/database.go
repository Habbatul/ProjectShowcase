package config

import (
	"database/sql"
	"log"
	"sync"
)

var (
	db   *sql.DB
	mu   sync.Mutex
	once sync.Once
)

// ConnectDB mengembalikan instance tunggal dari koneksi database
func ConnectDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=testfirstgo_db password=mysecretpassword sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	})

	return db
}
