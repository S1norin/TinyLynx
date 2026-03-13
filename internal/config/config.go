package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPassword string
	DBPort string
}

func Load() Config {
	err := godotenv.Load("")
	if err != nil {
		log.Println("No .env file")
	
	}
	
	dbUser := os.Getenv(("DB_USER"))
	if err != nil {
		dbUser = "postgres"
	}
	
	dbPassword := os.Getenv(("DB_PASSWORD"))
	if err != nil {
		dbPassword = "postgres"
	}
	
	dbPort := "5433"
	
	
	loadedConfig := Config{
		DBUser: dbUser,
		DBPassword: dbPassword,
		DBPort: dbPort,
	}
	
	return loadedConfig
}