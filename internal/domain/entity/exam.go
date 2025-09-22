package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ExamStatus string

const (
	StatusPending    ExamStatus = "pending"
	StatusProcessing ExamStatus = "processing"
	StatusDone       ExamStatus = "done"
	StatusFailed     ExamStatus = "failed"
)

type Exam struct {
	ID        string
	PatientID string
	ExamType  string
	Status    ExamStatus
	CreatedAt time.Time
}

func NewExam(patientID, examType string) *Exam {
	return &Exam{
		ID:        fmt.Sprintf("E-%s", uuid.New().String()), // Gera um ID único
		PatientID: patientID,
		ExamType:  examType,
		Status:    StatusPending, // Todo novo exame começa como 'pending'
		CreatedAt: time.Now(),
	}
}