package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser         string
	DBPassword     string
	DBPortExposed  string
	DBName         string
}

func Load() Config {
	err := godotenv.Load("")
	if err != nil {
		log.Println("No .env file")

	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dBPortExposed := os.Getenv("DB_PORT_EXPOSED")
	if dBPortExposed == "" {
		dBPortExposed = "5433"
	}

	dBName := os.Getenv("DB_NAME")
	if dBName == "" {
		dBName = "tinylynx"
	}

	loadedConfig := Config{
		DBUser:         dbUser,
		DBPassword:     dbPassword,
		DBPortExposed:  dBPortExposed,
		DBName:         dBName,
	}

	return loadedConfig
}
