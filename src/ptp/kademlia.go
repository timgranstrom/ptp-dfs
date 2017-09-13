package ptp

import (
	"github.com/golang/protobuf/proto"
	//"ptp/proto"
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
		print(contact.Address) //Just temporary to avoid errors when compiling
			// TODO Create and send request to contact
		}()
	}
	for {
		select {
			case reply := <- kademlia.replies: //Idle wait for replies to requests
				print(reply.String()) //Just temporary to avoid errors when compiling

				// TODO If reply has closer contacts than any in the contact list
					// TODO Push the new closer contact, pop the furthest
					// TODO Create and send request to new contact (in a routine)
			// TODO Timeout case for when requestees don't reply fast enough (might be disconnected/dead/slow)
		}
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
