package main

import (
	"exam-processing-service/internal/infra/database"
	"exam-processing-service/internal/infra/http/router"
	"log"
)

func main() {
	database.RunMigrations()
	r := router.SetupRouter()

	log.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
