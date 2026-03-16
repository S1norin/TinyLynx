package main

import (
	"tinylynx/internal/storage"
)

func main() {
	storage.RunMigrations()
}