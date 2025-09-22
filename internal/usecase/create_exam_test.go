package usecase

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"exam-processing-service/internal/usecase/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateExamUseCase_Execute(t *testing.T) {
	mockRepo := new(repository.ExamRepositoryMock)
	jobQueue := make(chan *entity.Exam, 1) // Usamos um channel real para o teste

	mockRepo.On("Save", mock.AnythingOfType("*entity.Exam")).Return(nil)

	uc := NewCreateExamUseCase(mockRepo, jobQueue)
	input := dto.CreateExamInputDTO{
		PatientID: "patient-123",
		ExamType:  "blood-test",
	}

	output, err := uc.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ExamID)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "Save", 1)

	assert.Len(t, jobQueue, 1)
	examInQueue := <-jobQueue
	assert.Equal(t, input.PatientID, examInQueue.PatientID)
	assert.Equal(t, entity.StatusPending, examInQueue.Status)
}
