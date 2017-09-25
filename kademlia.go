package ptp

import (
	//"github.com/golang/protobuf/proto"
	"container/list"
	//"ptp/proto"
	"log"
)

type Kademlia struct {
	routingTable RoutingTable
	network *Network
	workers *list.List //List of all current workers
	idCount int64 //Global id counter for workers
	boostrapNode *Contact //boostrap node to get you into the network
}

//Create a new kademlia with random id and at a specific address
func NewKademlia (address string, bootstrapNode *Contact) *Kademlia{
	myKadID := NewRandomKademliaID() //Create a new random kademlia id
	meContact := NewContact(myKadID,address)
	routingTable := *NewRoutingTable(meContact)
	kademlia := &Kademlia{
		routingTable:routingTable, //Create routing table with myself as contact
		network: NewNetwork(routingTable), //Create a new network
		workers: list.New(), //Create a new linked list for workers
		idCount:0,
	}
	if bootstrapNode != nil {
		kademlia.routingTable.AddContact(*bootstrapNode) //Add boostrap node in network
		log.Println(kademlia.routingTable.me.Address+": Added boostrap node as contact")
	}
	return kademlia
}

func (kademlia *Kademlia) Run(){
	go kademlia.network.Listen() //Start listener on network
	dispatcher := NewDispatcher(kademlia.network)
	go dispatcher.StartDispatcher() //always run
	kademlia.BoostrapProcess()
}

func (kademlia *Kademlia) BoostrapProcess(){
	go kademlia.LookupContact(&kademlia.routingTable.me) //Find closest nodes to me in network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	log.Println(kademlia.routingTable.me.Address+" :Lookup Contact was called internally")

	workRecievedCount := 0
	expectedWorkCount := 0

	worker := kademlia.NewWorker() //Create a worker that maps to the function
	kademlia.network.WorkerQueue <- worker //add the worker to the worker queue so that we can recieve messages
	contacts := ContactCandidates{} //closes contact candidates for the lookup
	contacts.Append(kademlia.routingTable.FindClosestContacts(target.ID,3)) //Retrieve nodes own closest contacts

	for _, contact := range contacts.contacts { //Send request to nodes own closests contacts for their closests contacts
		//go func() { //Send requests concurrently
		//print(contact.Address) //Just temporary to avoid errors when compiling

			// TODO Create and send lookupcontact request to contact
		//}()
		go kademlia.network.SendFindContactMessage(&contact,worker.id) //Create and send lookupcontact request to contact
	}
	for {
		select {
			case reply := <- worker.workRequest: //Idle wait for replies to requests
				workRecievedCount++ //increment received work counter when received reply
				if workRecievedCount < expectedWorkCount{
					kademlia.network.WorkerQueue <- worker //If we still expect more answers, re-add worker to queue
				}
				replyContactsProto := reply.message.GetMsg_2().Contacts //Get the Contact Proto message from the reply
				//Convert protbuf contact to kademlia contact
				replyContacts := []Contact{}
				for _, replyContact := range replyContactsProto{
					replyKademlia := NewKademliaID(*replyContact.KademliaId) //Create kademliaId from reply
					replyContact := NewContact(replyKademlia,*replyContact.Address) //Create contact from reply
					replyContacts = append(replyContacts, replyContact) //Append reply contacts into a list
				}
				addedContacts := contacts.AppendClosestContacts(replyContacts,3) //Add replied contacts to the contact candidate list

				for _,contact := range addedContacts{
					go kademlia.network.SendFindContactMessage(&contact,worker.id) //Fire off a new find contact request for each contact
				}
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
			case reply := <- worker.workRequest: //Idle wait for replies to requests
				print(reply.id) //Just temporary to avoid errors when compiling
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
			case reply := <- worker.workRequest: //Idle wait for replies to requests
				print(reply.id) //Just temporary to avoid errors when compiling
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
