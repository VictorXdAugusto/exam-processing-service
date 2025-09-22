// internal/usecase/get_exam_test.go
package usecase

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"exam-processing-service/internal/usecase/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetExamUseCase_Execute(t *testing.T) {
	mockRepo := new(repository.ExamRepositoryMock)
	expectedExam := &entity.Exam{
		ID:        "exam-123",
		PatientID: "patient-456",
		ExamType:  "x-ray",
		Status:    entity.StatusDone,
		CreatedAt: time.Now(),
	}

	mockRepo.On("FindByID", "exam-123").Return(expectedExam, nil)

	uc := NewGetExamUseCase(mockRepo)
	input := dto.GetExamInputDTO{ID: "exam-123"}

	output, err := uc.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByID", 1)

	assert.Equal(t, expectedExam.ID, output.ID)
	assert.Equal(t, expectedExam.PatientID, output.PatientID)
	assert.Equal(t, string(expectedExam.Status), output.Status)
}
