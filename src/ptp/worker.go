package ptp

import "github.com/golang/protobuf/proto"

type Worker struct {
	replies chan proto.Message //Functions own channels to receive messages in
	id int //Functions own id to attach in requests so replies can come back to the function through the dispatcher
}


func (kademlia *Kademlia) NewWorker() Worker{
	id := kademlia.idCount
	worker := Worker {
		replies: make(chan proto.Message),
		id: id,
	}
	kademlia.idCount++
	return worker
}
