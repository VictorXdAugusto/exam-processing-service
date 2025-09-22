package worker

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"log"
	"time"
)

type Worker struct {
	ID             int
	ExamRepository repository.ExamRepository
}

func NewWorker(id int, repo repository.ExamRepository) *Worker {
	return &Worker{
		ID:             id,
		ExamRepository: repo,
	}
}

func (w *Worker) ProcessJobs(jobQueue <-chan *entity.Exam) {
	for exam := range jobQueue {
		time.Sleep(5 * time.Second)
		log.Printf("Worker %d: iniciando processamento do exame %s", w.ID, exam.ID)

		exam.Status = entity.StatusProcessing
		if err := w.ExamRepository.Update(exam); err != nil {
			log.Printf("Worker %d: ERRO ao atualizar status para 'processing' do exame %s: %v", w.ID, exam.ID, err)
			continue
		}

		time.Sleep(5 * time.Second)

		if time.Now().UnixNano()%10 == 0 {
			exam.Status = entity.StatusFailed
		} else {
			exam.Status = entity.StatusDone
		}

		if err := w.ExamRepository.Update(exam); err != nil {
			log.Printf("Worker %d: ERRO ao atualizar status final do exame %s: %v", w.ID, exam.ID, err)
			continue
		}

		log.Printf("Worker %d: finalizou o processamento do exame %s com status %s", w.ID, exam.ID, exam.Status)
	}
}
