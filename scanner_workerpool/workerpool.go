package workerpool

import (
	"log"
)

var WorkerChannel = make(chan chan Work)

type Collector struct {
	Work chan Work
	End chan bool
}

func StartDispatcher(workerCount int, channelBuffer int, workFn DoWork) Collector {
	var i int
	var workers []Worker

	input := make(chan Work, channelBuffer)
	end := make(chan bool)
	collector := Collector{Work:input, End:end}

	for i < workerCount {
		i++

		log.Printf("spawning worker")

		worker := Worker{
			ID: i,
			Channel: make(chan Work),
			WorkerChannel: WorkerChannel,
			End: make(chan bool),
		}

		worker.Start(workFn)
		workers = append(workers, worker)
	}

	go func() {
		for {
			select {
			case <- end:
				for _, w := range workers {
					w.Stop()
				}
				return
			case work := <-input:
				worker := <-WorkerChannel
				worker <- work
			}
		}
	}()

	return collector
}
