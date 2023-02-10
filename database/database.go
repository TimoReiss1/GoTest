package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=1 dbname=product_db sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
}

func Close() {
	db.Close()
}
