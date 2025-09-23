package app

import (
	"context"
	"exam-processing-service/internal/config"
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/infra/database"
	"exam-processing-service/internal/infra/http/router"
	"exam-processing-service/internal/infra/repository"
	"exam-processing-service/internal/infra/worker"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	config     *config.Config
	server     *http.Server
	dbPool     *pgxpool.Pool
	jobQueue   chan *entity.Exam
	workerPool *WorkerPool
	logger     *log.Logger
}

type WorkerPool struct {
	workers []*worker.Worker
	wg      *sync.WaitGroup
}

func New(cfg *config.Config) *Application {
	return &Application{
		config: cfg,
		logger: log.New(os.Stdout, "[APP] ", log.LstdFlags|log.Lshortfile),
	}
}

func (a *Application) Initialize() error {
	a.logger.Println("Initializing application...")

	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	a.jobQueue = make(chan *entity.Exam, a.config.Worker.QueueSize)

	if err := a.initWorkerPool(); err != nil {
		return fmt.Errorf("failed to initialize worker pool: %w", err)
	}

	if err := a.initHTTPServer(); err != nil {
		return fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	a.logger.Println("Application initialized successfully")
	return nil
}

func (a *Application) Run() error {
	a.logger.Println("Starting application...")

	a.startWorkerPool()

	go func() {
		a.logger.Printf("Starting HTTP server on port %s...", a.config.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Printf("HTTP server error: %v", err)
		}
	}()

	a.waitForShutdown()

	a.logger.Println("Shutdown signal received, returning from Run()")
	return nil
}

func (a *Application) Shutdown() error {
	a.logger.Println("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Printf("HTTP server shutdown error: %v", err)
		return err
	}
	a.logger.Println("HTTP server stopped")

	a.shutdownWorkerPool()

	if a.dbPool != nil {
		a.dbPool.Close()
		a.logger.Println("Database connection closed")
	}

	a.logger.Println("Application shutdown completed")
	return nil
}

func (a *Application) initDatabase() error {
	a.logger.Println("Running database migrations...")
	database.RunMigrations()

	a.logger.Println("Connecting to database...")
	dbPool, err := database.NewPostgresConnection()
	if err != nil {
		return err
	}

	a.dbPool = dbPool
	a.logger.Println("Database connected successfully")
	return nil
}

func (a *Application) initWorkerPool() error {
	a.logger.Printf("Initializing worker pool with %d workers...", a.config.Worker.Count)

	examRepo := repository.NewExamRepositoryPostgres(a.dbPool)
	wg := &sync.WaitGroup{}
	workers := make([]*worker.Worker, a.config.Worker.Count)

	for i := 0; i < a.config.Worker.Count; i++ {
		workers[i] = worker.NewWorker(i+1, examRepo, wg)
	}

	a.workerPool = &WorkerPool{
		workers: workers,
		wg:      wg,
	}

	a.logger.Printf("Worker pool initialized with %d workers", a.config.Worker.Count)
	return nil
}

func (a *Application) initHTTPServer() error {
	a.logger.Println("Initializing HTTP server...")

	examRepo := repository.NewExamRepositoryPostgres(a.dbPool)
	r := router.SetupRouter(examRepo, a.jobQueue)

	a.server = &http.Server{
		Addr:    ":" + a.config.Server.Port,
		Handler: r,
	}

	a.logger.Println("HTTP server initialized")
	return nil
}

func (a *Application) startWorkerPool() {
	a.logger.Println("Starting worker pool...")

	for i, w := range a.workerPool.workers {
		go w.ProcessJobs(a.jobQueue)
		a.logger.Printf("Worker %d started", i+1)
	}

	a.logger.Printf("All %d workers started", len(a.workerPool.workers))
}

func (a *Application) shutdownWorkerPool() {
	a.logger.Println("Shutting down worker pool...")

	close(a.jobQueue)
	a.logger.Println("Job queue closed. Waiting for workers to finish...")

	a.workerPool.wg.Wait()
	a.logger.Println("All workers finished")
}

func (a *Application) waitForShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	a.logger.Printf("Received shutdown signal: %v", sig)
}
