package ptp

import (
	"fmt"
	//"github.com/golang/protobuf/proto"
	//"ptp/proto"
	"log"
	"strconv"
)

var WorkerQueue chan Worker

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan Worker, nworkers)

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
				log.Println("Received work requeust")
				go func() { //New GoRoutine
					for{ //Always try to find the correct worker until Timeout or success
						worker := <- WorkerQueue
						if worker.id == work.id { //If worker and work id match, SUCCESS!
							log.Println("Dispatching work request for "+strconv.Itoa(int(work.id)))
							worker.workRequest <- work //Dispatch to the correct worker instead (not correct now)
							//WorkerQueue <- worker
							break
						} else {
							WorkerQueue <- worker
						}
					}

				}()
			}
		}
	}()
}
