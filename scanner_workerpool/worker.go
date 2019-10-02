package workerpool

import (
	"log"
)

type DoWork func(id string, input string)

type Work struct {
	ID string
	Job string
}

type Worker struct {
	ID int
	WorkerChannel chan chan Work
	Channel chan Work
	End chan bool
}

func (w *Worker) Start(fn DoWork) {
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				log.Printf("Running job on worker[%d] = %s", w.ID, job.ID)
				fn(job.ID, job.Job)
			case <- w.End:
				log.Printf("Killing worker[%d]", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.End <- true
}
