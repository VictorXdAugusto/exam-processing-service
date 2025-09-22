package handler

import (
	"exam-processing-service/internal/usecase"
	"exam-processing-service/internal/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	CreateExamUseCase *usecase.CreateExamUseCase
}

func NewExamHandler(createExamUC *usecase.CreateExamUseCase) *ExamHandler {
	return &ExamHandler{
		CreateExamUseCase: createExamUC,
	}
}

func (h *ExamHandler) CreateExam(c *gin.Context) {
	var input dto.CreateExamInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json provided"})
		return
	}

	output, err := h.CreateExamUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create exam"})
		return
	}

	c.JSON(http.StatusCreated, output)
}
