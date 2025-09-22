package dto

import "time"

type CreateExamInputDTO struct {
	PatientID string `json:"patient_id"`
	ExamType  string `json:"exam_type"`
}

type CreateExamOutputDTO struct {
	ExamID string `json:"exam_id"`
}

type GetExamInputDTO struct {
	ID string
}

type GetExamOutputDTO struct {
	ID        string    `json:"exam_id"`
	PatientID string    `json:"patient_id"`
	ExamType  string    `json:"exam_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
