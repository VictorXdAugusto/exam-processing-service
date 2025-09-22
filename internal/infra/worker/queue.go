package worker

import "exam-processing-service/internal/domain/entity"

// Podemos colocar até 100 exames na fila antes que o produtor precise esperar.
var JobQueue = make(chan *entity.Exam, 100)
