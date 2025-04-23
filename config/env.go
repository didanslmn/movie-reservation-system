package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("errorload .env file")
	}
	log.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))
}
