package ptp

import (
	"log"
	"strconv"
	"github.com/timgranstrom/ptp-dfs/protoMessages"
)

type Dispatcher struct{
	network *Network
}

func NewDispatcher(network *Network) *Dispatcher{
	dispatch := &Dispatcher{network:network}
	return dispatch
}



func (dispatcher *Dispatcher) StartDispatcher() {
	// First, initialize the channel we are going to but the workers' work channels into.
	//WorkerQueue = make(chan Worker, nworkers)

	// Now, create all of our workers. NOT NECESSARY
	/*for i := 0; i<nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}*/

	go func() { //New GoRoutine
		for {
			select {
			case work := <-dispatcher.network.WorkQueue:
				log.Println(dispatcher.network.routingTable.me.Address,": Received work request, is reply:",*work.message.IsReply)
				go func() { //New GoRoutine
					if *work.message.IsReply {
						for { //Always try to find the correct worker until Timeout or success
							worker := <-dispatcher.network.WorkerQueue
							if worker.id == work.id { //If worker and work id match, SUCCESS!
								log.Println(dispatcher.network.routingTable.me.Address,": Dispatching work request for " + strconv.Itoa(int(work.id)))
								worker.workRequest <- work //Dispatch to the correct worker instead (not correct now)
								break
							} else {
								dispatcher.network.WorkerQueue <- worker //Add worker back to queue if it was the wrong worker
							}
						}
					} else{
						switch *work.message.MessageType{
						case protoMessages.MessageType_PING:
							log.Println(dispatcher.network.routingTable.me.Address," :Picked up PING from Message Queue")
							break
						case protoMessages.MessageType_FIND_CONTACT:
							log.Println(dispatcher.network.routingTable.me.Address," :Picked up FIND CONTACT from Message Queue")
							dispatcher.network.RecieveFindContactMessage(&work)
							break
						case protoMessages.MessageType_FIND_DATA:
							log.Println(dispatcher.network.routingTable.me.Address," :Picked up FIND DATA from Message Queue")
							break
						case protoMessages.MessageType_SEND_STORE:
							log.Println(dispatcher.network.routingTable.me.Address," :Picked up SEND STORE from Message Queue")
							break
						default:
							log.Println(dispatcher.network.routingTable.me.Address," :Picked up UNKNOWN MESSAGE TYPE from Message Queue")
						}
					}
				}()
			}
		}
	}()
}
