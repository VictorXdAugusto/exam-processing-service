package usecase

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"exam-processing-service/internal/usecase/dto"
)

type CreateExamUseCase struct {
	ExamRepository repository.ExamRepository
	JobQueue       chan<- *entity.Exam
}

func NewCreateExamUseCase(repo repository.ExamRepository, jobQueue chan<- *entity.Exam) *CreateExamUseCase {
	return &CreateExamUseCase{
		ExamRepository: repo,
		JobQueue:       jobQueue,
	}
}

func (uc *CreateExamUseCase) Execute(input dto.CreateExamInputDTO) (*dto.CreateExamOutputDTO, error) {
	exam := entity.NewExam(input.PatientID, input.ExamType)
	err := uc.ExamRepository.Save(exam)
	if err != nil {
		return nil, err
	}

	uc.JobQueue <- exam

	output := &dto.CreateExamOutputDTO{
		ExamID: exam.ID,
	}
	return output, nil
}
