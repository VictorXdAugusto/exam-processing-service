package repository

import (
	"context"
	"exam-processing-service/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExamRepositoryPostgres struct {
	DB *pgxpool.Pool
}

func NewExamRepositoryPostgres(db *pgxpool.Pool) *ExamRepositoryPostgres {
	return &ExamRepositoryPostgres{DB: db}
}

func (r *ExamRepositoryPostgres) Save(exam *entity.Exam) error {
	sql := `INSERT INTO exams (id, patient_id, exam_type, status, created_at) 
	         VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(context.Background(), sql,
		exam.ID,
		exam.PatientID,
		exam.ExamType,
		exam.Status,
		exam.CreatedAt,
	)
	return err
}

func (r *ExamRepositoryPostgres) Update(exam *entity.Exam) error {
	sql := `UPDATE exams SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(context.Background(), sql, exam.Status, exam.ID)
	return err
}

func (r *ExamRepositoryPostgres) FindByID(id string) (*entity.Exam, error) {
	sql := `SELECT id, patient_id, exam_type, status, created_at 
	         FROM exams WHERE id = $1`

	row := r.DB.QueryRow(context.Background(), sql, id)

	var exam entity.Exam
	err := row.Scan(
		&exam.ID,
		&exam.PatientID,
		&exam.ExamType,
		&exam.Status,
		&exam.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &exam, nil
}
