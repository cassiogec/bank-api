package main

import (
	"log"
	"os"

	"github.com/cassiogec/bank-api/controllers"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":8080")
}
