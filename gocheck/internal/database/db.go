package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "user=root password=password dbname=gocheck_db sslmode=disable host=localhost port=5432"
	var err error
	
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Nie udało się otworzyć bazy: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Baza nie odpowiada: ", err)
	}

	log.Println("Połączono z bazą PostgreSQL!")
}