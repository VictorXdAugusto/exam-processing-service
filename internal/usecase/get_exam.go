package usecase

import (
	"exam-processing-service/internal/domain/repository"
	"exam-processing-service/internal/usecase/dto"
)

type GetExamUseCase struct {
	ExamRepository repository.ExamRepository
}

func NewGetExamUseCase(repo repository.ExamRepository) *GetExamUseCase {
	return &GetExamUseCase{
		ExamRepository: repo,
	}
}

func (uc *GetExamUseCase) Execute(input dto.GetExamInputDTO) (*dto.GetExamOutputDTO, error) {
	exam, err := uc.ExamRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	output := &dto.GetExamOutputDTO{
		ID:        exam.ID,
		PatientID: exam.PatientID,
		ExamType:  exam.ExamType,
		Status:    string(exam.Status),
		CreatedAt: exam.CreatedAt,
	}

	return output, nil
}
