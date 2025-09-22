package main

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/infra/database"
	"exam-processing-service/internal/infra/http/router"
	"exam-processing-service/internal/infra/repository"
	"exam-processing-service/internal/infra/worker"
	"log"
)

func main() {
	database.RunMigrations()

	dbPool, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}
	defer dbPool.Close()

	examRepo := repository.NewExamRepositoryPostgres(dbPool)

	log.Println("Starting workers...")
	const numberOfWorkers = 5
	jobQueue := make(chan *entity.Exam, 100)
	for i := 1; i <= numberOfWorkers; i++ {
		w := worker.NewWorker(i, examRepo)
		go w.ProcessJobs(jobQueue)
	}
	log.Printf("%d workers started.", numberOfWorkers)

	r := router.SetupRouter(examRepo, jobQueue)

	log.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
