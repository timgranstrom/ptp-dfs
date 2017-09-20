package ptp

import (
	"fmt"
	//"github.com/golang/protobuf/proto"
	//"ptp/proto"
)

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers. NOT NECESSARY
	/*for i := 0; i<nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}*/

	go func() { //New GoRoutine
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work requeust")
				go func() { //New GoRoutine
					worker := <-WorkerQueue

					fmt.Println("Dispatching work request")
					worker <- work //Dispatch to the correct worker instead (not correct now)
				}()
			}
		}
	}()
}
