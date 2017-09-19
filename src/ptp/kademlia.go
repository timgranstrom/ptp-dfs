package ptp

import (
	//"github.com/golang/protobuf/proto"
	"container/list"
)

type Kademlia struct {
	routingTable RoutingTable
	network *Network
	workers *list.List //List of all current workers
	idCount int //Global id counter for workers
}

func NewKademlia () *Kademlia{
	myKadID := NewRandomKademliaID() //Create a new random kademlia id
	meContact := NewContact(myKadID,"localhost") //For now use localhost as example for IP address
	kademlia := &Kademlia{
		routingTable:*NewRoutingTable(meContact), //Create routing table with myself as contact
		network: NewNetwork(), //Create a new network
		workers: list.New(), //Create a new linked list for workers
		idCount:0,
	}
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	worker := kademlia.NewWorker()
	contacts := kademlia.routingTable.FindClosestContacts(target.ID,3) //Retrieve nodes own closest contacts
	for _, contact := range contacts { //Send request to nodes own closests contacts for their closests contacts
		go func() { //Send requests concurrently
		print(contact.Address) //Just temporary to avoid errors when compiling
			// TODO Create and send lookupcontact request to contact
		}()
	}
	for {
		select {
			case reply := <- worker.replies: //Idle wait for replies to requests
				print(reply.String()) //Just temporary to avoid errors when compiling

				// TODO If reply has closer contacts than any in the contact list
					// TODO Push the new closer contact, pop the furthest
					// TODO If it was last reply
						// TODO Close the loop
					// TODO Goroutine to create and send request to new contact
			// TODO Timeout case for when requestees don't reply fast enough (might be disconnected/dead/slow)
				// TODO Close the loop
		}
	}
	// TODO Return closest contacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	worker := kademlia.NewWorker()
	contacts := kademlia.routingTable.FindClosestContacts(NewKademliaID(hash),3) //Retrieve nodes own closest contacts
	for _, contact := range contacts { //Send request to nodes own closests contacts for their closests contacts
		go func() { //Send requests concurrently
			print(contact.String()) //Just temporary to avoid errors when compiling
			// TODO Create and send lookupdata request to contact
		}()
	}
	for {
		select {
			case reply := <- worker.replies: //Idle wait for replies to requests
				print(reply.String()) //Just temporary to avoid errors when compiling
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
	worker := kademlia.NewWorker()
	contacts := kademlia.routingTable.FindClosestContacts(NewKademliaID(string(data)),3) //Retrieve nodes own closest contacts
	for _, contact := range contacts { //Send request to nodes own closests contacts for their closests contacts
		go func() { //Send requests concurrently
			print(contact.Address) //Just temporary to avoid errors when compiling
			// TODO Create and send lookupcontact request to contact
		}()
	}
	for {
		select {
			case reply := <- worker.replies: //Idle wait for replies to requests
				print(reply.String()) //Just temporary to avoid errors when compiling
				// TODO If reply has closer contacts than any in the contact list
					// TODO Push the new closer contact, pop the furthest
					// TODO Goroutine to create and send request to new contact
				// TODO Else...
					// TODO If it was last reply
						// TODO Close the loop
			// TODO Timeout case for when requestees don't reply fast enough (might be disconnected/dead/slow)
				// TODO Close the loop
		}
	}
	//Now we've found closest nodes to send store requests to
	for _, contact := range contacts {
		go func() {
			print(contact.Address) //Just temporary to avoid errors when compiling
			// TODO Create and send store request to contact
		}()
	}
}
