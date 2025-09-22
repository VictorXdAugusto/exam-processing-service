package dto

type CreateExamInputDTO struct {
	PatientID string `json:"patient_id"`
	ExamType  string `json:"exam_type"`
}

type CreateExamOutputDTO struct {
	ExamID string `json:"exam_id"`
}
