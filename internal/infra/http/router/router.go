package router

import (
	"exam-processing-service/internal/infra/database"
	"exam-processing-service/internal/infra/http/handler"
	"exam-processing-service/internal/infra/repository"
	"exam-processing-service/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	dbPool, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	examRepo := repository.NewExamRepositoryPostgres(dbPool)
	createExamUseCase := usecase.NewCreateExamUseCase(examRepo)

	healthHandler := handler.NewHealthHandler()
	examHandler := handler.NewExamHandler(createExamUseCase)

	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler.HealthCheck)
		api.POST("/exams", examHandler.CreateExam)
	}

	return router
}
