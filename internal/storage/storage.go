package storage

import (
	"context"
	"fmt"
	"log"
	"sync"
	"database/sql"
	"tinylynx/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var conn *pgx.Conn
var pool *pgxpool.Pool
var oncePool, onceConn sync.Once

func getDsn() string {
	cfg := config.Load()
	return fmt.Sprintf("postgres://%s:%s@db:5432/%s", cfg.DBUser, cfg.DBPassword, "tinylynx")
}

func GetConn(ctx context.Context) *pgx.Conn {
	oncePool.Do(
		func() {
			dbConn, err := pgx.Connect(ctx, getDsn())
			if err != nil {
				log.Fatalf("storage.go getConn(): Couldn't get connection %v", err)
			}
			conn = dbConn
		})
	log.Print("storage.go getConn(): Got connection")
	return conn
}

func GetPool(ctx context.Context) *pgxpool.Pool {
	oncePool.Do(
		func() {
			dbPool, err := pgxpool.New(ctx, getDsn())
			if err != nil {
				log.Fatalf("storage.go getPool(): Couldn't get pool %v", err)
			}
			pool = dbPool
		})
	log.Print("storage.go getPool(): Got pool")
	return pool
}

func RunMigrations(ctx context.Context) {
	db, err := sql.Open("pgx", getDsn())
	if err != nil {
		log.Fatalf("storage.go RunMigrations(): Couldn't get connection %v", err)
	}

	defer db.Close()

	err = goose.Up(db, "internal/storage")
	if err != nil {
		log.Fatalf("storage.go RunMigrations(): Couldn't go up with migration %v", err)
	}
	
}