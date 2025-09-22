package repository

import (
	"exam-processing-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type ExamRepositoryMock struct {
	mock.Mock
}

func (m *ExamRepositoryMock) Save(exam *entity.Exam) error {
	args := m.Called(exam)
	return args.Error(0)
}

func (m *ExamRepositoryMock) Update(exam *entity.Exam) error {
	args := m.Called(exam)
	return args.Error(0)
}

func (m *ExamRepositoryMock) FindByID(id string) (*entity.Exam, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Exam), args.Error(1)
}
