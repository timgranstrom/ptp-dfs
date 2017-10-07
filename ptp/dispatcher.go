package ptp

import (
	"log"
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
	log.Println(dispatcher.network.routingTable.me.Address + ": Dispatcher active")

	go func() {
		for {
			select {
				//Wait for messages that are either replies or requests
				case workReceived := <-dispatcher.network.WorkQueue:
					go func(work WorkRequest) {
						//If it's a reply, send it to a worker, otherwise it's a request
						if work.message.IsReply {
							for { //Always try to find the correct worker until Timeout or success
								worker := <-dispatcher.network.WorkerQueue
								//See if the worker is active, otherwise close it and don't put it back
								if worker.active {
									//See if it is the right worker to send the work to, otherwise put it pack
									if worker.id == work.id {
										worker.workRequest <- work //Dispatch to the correct worker
										break
									} else {
										dispatcher.network.WorkerQueue <- worker //Put the worker back
									}
								} else {
									close(worker.workRequest)
								}
							}
						} else {
							switch work.message.MessageType {
							case protoMessages.MessageType_PING:
								dispatcher.network.ReceivePingContactMessage(&work)
								break
							case protoMessages.MessageType_FIND_CONTACT:
								dispatcher.network.RecieveFindContactMessage(&work)
								break
							case protoMessages.MessageType_FIND_DATA:
								dispatcher.network.ReceiveFindDataMessage(&work)
								break
							case protoMessages.MessageType_SEND_STORE:
								dispatcher.network.RecieveStoreMessage(&work)
								break
							default:
								break
							}
						}
					}(workReceived)
			}
		}
	}()
}
