package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vjurczenia/mblg/mblg"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	mblg.StartServer()
}
