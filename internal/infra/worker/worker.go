package worker

import (
	"exam-processing-service/internal/domain/entity"
	"exam-processing-service/internal/domain/repository"
	"log"
	"sync"
	"time"
)

type Worker struct {
	ID             int
	ExamRepository repository.ExamRepository
	Wg             *sync.WaitGroup
}

func NewWorker(id int, repo repository.ExamRepository, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:             id,
		ExamRepository: repo,
		Wg:             wg,
	}
}

func (w *Worker) ProcessJobs(jobQueue <-chan *entity.Exam) {
	for exam := range jobQueue {
		w.Wg.Add(1)

		go func(e *entity.Exam) {
			defer w.Wg.Done()

			log.Printf("Worker %d: iniciando processamento do exame %s", w.ID, e.ID)

			e.Status = entity.StatusProcessing
			if err := w.ExamRepository.Update(e); err != nil {
				log.Printf("Worker %d: ERRO ao atualizar status para 'processing' do exame %s: %v", w.ID, e.ID, err)
				return
			}

			time.Sleep(5 * time.Second)

			if time.Now().UnixNano()%10 == 0 {
				e.Status = entity.StatusFailed
			} else {
				e.Status = entity.StatusDone
			}

			if err := w.ExamRepository.Update(e); err != nil {
				log.Printf("Worker %d: ERRO ao atualizar status final do exame %s: %v", w.ID, e.ID, err)
				return
			}

			log.Printf("Worker %d: finalizou o processamento do exame %s com status %s", w.ID, e.ID, e.Status)
		}(exam)
	}
}
