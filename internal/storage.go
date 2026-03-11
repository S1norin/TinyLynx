package internal

import (
	"database/sql"
	"fmt"
	"storage"
	
	"github.com/jackc/pgx/v5"
)

func connection() {
	
}

func initilizeTable() {
	db, err := sql.Open("pgx", "")
}