package handler

import (
	"exam-processing-service/internal/usecase"
	"exam-processing-service/internal/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	CreateExamUseCase *usecase.CreateExamUseCase
	GetExamUseCase    *usecase.GetExamUseCase
}

func NewExamHandler(createExamUC *usecase.CreateExamUseCase, getExamUC *usecase.GetExamUseCase) *ExamHandler {
	return &ExamHandler{
		CreateExamUseCase: createExamUC,
		GetExamUseCase:    getExamUC,
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

func (h *ExamHandler) GetExam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exam id is required"})
		return
	}

	input := dto.GetExamInputDTO{ID: id}

	output, err := h.GetExamUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exam not found"})
		return
	}

	c.JSON(http.StatusOK, output)
}
