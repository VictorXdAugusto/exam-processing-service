package router

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"exam-processing-service/internal/infra/http/handler"
	"exam-processing-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	examRepo repository.ExamRepository,
	jobQueue chan<- *entity.Exam,
) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	createExamUseCase := usecase.NewCreateExamUseCase(examRepo, jobQueue)
	getExamUseCase := usecase.NewGetExamUseCase(examRepo)

	healthHandler := handler.NewHealthHandler()
	examHandler := handler.NewExamHandler(createExamUseCase, getExamUseCase)

	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler.HealthCheck)
		api.POST("/exams", examHandler.CreateExam)
		api.GET("/exams/:id", examHandler.GetExam)
	}

	return router
}
