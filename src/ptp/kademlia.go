package ptp

import (
	"github.com/golang/protobuf/proto"
	"ptp/proto"
)

type Kademlia struct {
	routingTable RoutingTable
	network Network
	replies chan proto.Message
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	contacts := kademlia.routingTable.FindClosestContacts(target.ID,3) //Retrieve nodes own closest contacts
	for _, contact := range contacts { //Send request to nodes own closests contacts for their closests contacts
		go func() { //Send requests concurrently
			// TODO Create and send request to contact
		}()
	}
	for {
		select {
			case reply := <- kademlia.replies: //Idle wait for replies to requests
				// TODO If reply has closer contacts than any in the contact list
					// TODO Push the new closer contact, pop the furthest
					// TODO If it was last reply
						// TODO Close the loop
					// TODO Goroutine to create and send request to new contact
			// TODO Timeout case for when requestees don't reply fast enough (might be disconnected/dead/slow)
				// TODO Close the loop
		}
	}
	// TODO Return list of closest contacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	contacts := kademlia.routingTable.FindClosestContacts(target.ID,3) //Retrieve nodes own closest contacts
	for _, contact := range contacts { //Send request to nodes own closests contacts for their closests contacts
		go func() { //Send requests concurrently
			// TODO Create and send request to contact
		}()
	}
	for {
		select {
			case reply := <- kademlia.replies: //Idle wait for replies to requests
				// TODO If reply has address of data
					// TODO Download data
					// TODO Create and send store request to next closest node
					// TODO Close the loop
				// TODO If reply has closer contacts than any in the contact list
					// TODO Push the new closer contact, pop the furthest
					// TODO Goroutine to create and send request to new contact
			// TODO Timeout case for when requestees don't reply fast enough (might be disconnected/dead/slow)
				// TODO Close the loop
		}
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO Goroutine to create and send store request for data
}
