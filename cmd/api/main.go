package main

import (
	"exam-processing-service/internal/app"
	"exam-processing-service/internal/config"
	"log"
)

func main() {
	cfg := config.Load()

	application := app.New(cfg)

	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application run error: %v", err)
	}

	if err := application.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown application gracefully: %v", err)
	}
}
