package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"tinylynx/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var pool *pgxpool.Pool
var oncePool, onceConn sync.Once

func getDsn() string {
	cfg := config.Load()
	return fmt.Sprintf("postgres://%s:%s@db:5432/%s", cfg.DBUser, cfg.DBPassword, "tinylynx")
}

func GetPool() *pgxpool.Pool {
	oncePool.Do(
		func() {
			dbPool, err := pgxpool.New(context.Background(), getDsn())
			if err != nil {
				log.Fatalf("storage.go getPool(): Couldn't get pool %v", err)
			}
			pool = dbPool
		})
	log.Print("storage.go getPool(): Got pool")
	return pool
}


func RunMigrations() {
	db, err := sql.Open("pgx", getDsn())
	if err != nil {
		log.Fatalf("storage.go RunMigrations(): Couldn't get connection %v", err)
	}

	defer db.Close()

	err = goose.Up(db, "internal/storage/migrations")
	if err != nil {
		log.Fatalf("storage.go RunMigrations(): Couldn't go up with migration %v", err)
	}

}