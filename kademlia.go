package ptp

import (
	//"github.com/golang/protobuf/proto"
	"container/list"
	//"ptp/proto"
	"log"
	"time"
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
		log.Println(kademlia.routingTable.me.Address+": Add boostrap node as contact")
		kademlia.routingTable.AddContact(*bootstrapNode) //Add boostrap node in network
	}
	return kademlia
}

func (kademlia *Kademlia) Run(){
	println("\n---------------")
	log.Println(kademlia.network.routingTable.me.Address,": NODE STARTED")
	go kademlia.network.Listen() //Start listener on network
	go kademlia.network.Sender() //Start the sender on network
	dispatcher := NewDispatcher(kademlia.network)
	dispatcher.StartDispatcher() //always run

	go kademlia.BoostrapProcess()

}

func (kademlia *Kademlia) BoostrapProcess(){
	for {
		if !kademlia.network.listenerActive {
			time.Sleep(time.Second / 10)
		} else{
			break
		}
	}
	log.Println(kademlia.network.routingTable.me.Address,"###########Bootstrapping process###########")
	kademlia.LookupContact(&kademlia.routingTable.me) //Find closest nodes to me in network
	log.Println(kademlia.network.routingTable.me.Address,"###########Finished Bootstrapping process###########")
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	workRecievedCount := 0
	expectedWorkCount := 0

	worker := kademlia.NewWorker() //Create a worker that maps to the function
	log.Println(kademlia.routingTable.me.Address+" :[Lookup Contact] was called internally, ID: ",worker.id)

	kademlia.network.WorkerQueue <- worker //add the worker to the worker queue so that we can recieve messages

	contacts := ContactCandidates{} //closes contact candidates for the lookup
	contacts.Append(kademlia.routingTable.FindClosestContacts(target.ID,3)) //Retrieve nodes own closest contacts

	log.Println(kademlia.routingTable.me.Address,":LIST OF CLOSEST CONTACTS INTERNALLY TO ",target.Address)

	for _, contact := range contacts.contacts { //Send request to nodes own closests contacts for their closests contacts
		log.Println(kademlia.routingTable.me.Address,":",contact.Address)
		expectedWorkCount++
		//TODO: Check if adding go below messes it up????
		kademlia.network.SendFindContactMessage(target,&contact,worker.id) //Create and send lookupcontact request to contact
	}
	log.Println(kademlia.routingTable.me.Address,":FINISHED LIST OF CLOSEST CONTACTS INTERNALLY")

lookForRepliesChannel:
		for {
			select {
				case reply := <- worker.workRequest: //Idle wait for replies to requests
					workRecievedCount++ //increment received work counter when received reply
					log.Println(kademlia.routingTable.me.Address,":--------------------- Work received count:",workRecievedCount)

					replyContactsProto := reply.message.GetMsg_2().Contacts //Get the Contact Proto message from the reply
					//Convert protbuf contact to kademlia contact
					replyContacts := []Contact{}
					for _, replyContact := range replyContactsProto{
						if *replyContact.KademliaId != kademlia.network.routingTable.me.ID.String(){
							replyKademlia := NewKademliaID(*replyContact.KademliaId) //Create kademliaId from reply
							replyContact := NewContact(replyKademlia,*replyContact.Address) //Create contact from reply
							replyContact.CalcDistance(target.ID) //Calculate distance
							replyContacts = append(replyContacts, replyContact) //Append reply contacts into a list
						}else{
							log.Println(kademlia.routingTable.me.Address,":Filtered from sending [Find Contact] message to self")
						}
					}

					println("CURRENT CONTACTS:")
					for _, elem := range contacts.contacts{
						println(elem.Address)
					}
					println("POSSIBLE NEW CONTACTS:")
					for _, elem := range replyContacts{
						println(elem.Address)
					}

					newContacts := contacts.AppendClosestContacts(replyContacts,3) //Add replied contacts to the contact candidate list

					println("NEW CONTACTS")
					for _,elem := range newContacts{
						println(elem.Address)
					}
					println("ALL CONTACTS")
					for _,elem := range contacts.contacts{
						println(elem.Address)
					}

					for _,contact := range newContacts{
						expectedWorkCount++
						//TODO: SET GO ROUTINE HERE AGAIN
						kademlia.network.SendFindContactMessage(target,&contact,worker.id) //Fire off a new find contact request for each contact
						}

					log.Println(kademlia.routingTable.me.Address,":---------------------  Work expected count:",expectedWorkCount)

					if workRecievedCount < expectedWorkCount{
						log.Println(kademlia.routingTable.me.Address,":!!!!!---------------------  ADDED worker back to queue, ID: ",worker.id)
						kademlia.network.WorkerQueue <- worker //If we still expect more answers, re-add worker to queue
					}
					// TODO If reply has closer contacts than any in the contact list
						// TODO Push the new closer contact, pop the furthest
						// TODO If it was last reply
							// TODO Close the loop
						// TODO Goroutine to create and send request to new contact
				// TODO Timeout case for when requestees don't reply fast enough (might be disconnected/dead/slow)
					// TODO Close the loop
			default:
				if workRecievedCount == expectedWorkCount{
					break lookForRepliesChannel
				}

			}
		}
	log.Println(kademlia.routingTable.me.Address,": Finished [Find Contact]\n")

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
