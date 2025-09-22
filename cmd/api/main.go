package main

import (
	"context"
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/infra/database"
	"exam-processing-service/internal/infra/http/router"
	"exam-processing-service/internal/infra/repository"
	"exam-processing-service/internal/infra/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	_, stop := context.WithCancel(context.Background())
	defer stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	database.RunMigrations()
	dbPool, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}
	defer dbPool.Close()

	examRepo := repository.NewExamRepositoryPostgres(dbPool)

	var wg sync.WaitGroup
	const numberOfWorkers = 5
	jobQueue := make(chan *entity.Exam, 100)

	log.Println("Starting workers...")
	for i := 1; i <= numberOfWorkers; i++ {
		w := worker.NewWorker(i, examRepo, &wg)
		go w.ProcessJobs(jobQueue)
	}
	log.Printf("%d workers started.", numberOfWorkers)

	r := router.SetupRouter(examRepo, jobQueue)
	server := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		log.Println("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	<-signalChan
	log.Println("Shutdown signal received. Shutting down gracefully...")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}
	log.Println("HTTP server stopped.")

	close(jobQueue)
	log.Println("Job queue closed. Waiting for workers to finish...")

	wg.Wait()
	log.Println("All workers have finished.")

	log.Println("Server gracefully stopped.")
}
