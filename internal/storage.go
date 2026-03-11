package internal

import (
	"database/sql"
	"fmt"
	"storage"
	
	"github.com/jackc/pgx/v5"
)

func connection() {
	dsn := ":password@tcp(localhost:3306)/dbname"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")
}



func initilizeTable() {
	db, err := sql.Open("pgx", "")
}