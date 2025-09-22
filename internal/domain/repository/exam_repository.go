package repository

import "exam-processing-service/internal/domain/entity"

type ExamRepository interface {
	Save(exam *entity.Exam) error
	Update(exam *entity.Exam) error
}
