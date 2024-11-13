package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
