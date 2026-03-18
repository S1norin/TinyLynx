package main

import (
	"fmt"
	"tinylynx/internal/storage"
)

func main() {
	storage.RunMigrations()
	fmt.Print()
}