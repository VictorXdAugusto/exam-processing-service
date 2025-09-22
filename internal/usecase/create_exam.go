package usecase

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"exam-processing-service/internal/usecase/dto"
)

type CreateExamUseCase struct {
	ExamRepository repository.ExamRepository
}

func NewCreateExamUseCase(repo repository.ExamRepository) *CreateExamUseCase {
	return &CreateExamUseCase{
		ExamRepository: repo,
	}
}

func (uc *CreateExamUseCase) Execute(input dto.CreateExamInputDTO) (*dto.CreateExamOutputDTO, error) {
	exam := entity.NewExam(input.PatientID, input.ExamType)

	err := uc.ExamRepository.Save(exam)
	if err != nil {
		return nil, err
	}

	output := &dto.CreateExamOutputDTO{
		ExamID: exam.ID,
	}
	return output, nil
}
