package ptp

import (
	//"fmt"
	"time"
	"github.com/timgranstrom/ptp-dfs/protoMessages"
)

type WorkRequest struct{
	id int64 //Functions own id to attach in requests so replies can come back to the function through the dispatcher
	message protoMessages.WrapperMessage //Wrapper message
}

type Worker struct {
	workRequest chan WorkRequest //Functions own channels to receive messages in (this is work)
	id int64 //Functions own id to attach in requests so replies can come back to the function through the dispatcher
	//WorkerQueue chan chan WorkRequest //Queue with workers
	timeoutDelay time.Duration //Time before the work times out
}

//Create a new worker and give it the Worker queue from the dispatcher
func (kademlia *Kademlia) NewWorker() Worker{
	id := kademlia.idCount
	worker := Worker {
		workRequest: make(chan WorkRequest),
		id: id,
		timeoutDelay: 10, //10 seconds before timeout on the entire function (stops listening for answer)
	}
	kademlia.idCount++
	return worker
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
/*
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.workRequest

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Received work request, delaying for %f seconds\n", w.ID, work.Delay.Seconds())

				time.Sleep(work.delay)
				fmt.Printf("worker%d: Hello, %s!\n", w.ID, work.Name)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}
*/
