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
	//wrcTest := make(chan int)

	worker := kademlia.NewWorker() //Create a worker that maps to the function
	//log.Println(kademlia.routingTable.me.Address+" :[Lookup Contact] was called internally, ID: ",worker.id)

	kademlia.network.WorkerQueue <- worker //add the worker to the worker queue so that we can recieve messages

	contacts := ContactCandidates{} //closes contact candidates for the lookup
	contacts.Append(kademlia.routingTable.FindClosestContacts(target.ID,3)) //Retrieve nodes own closest contacts

	//log.Println(kademlia.routingTable.me.Address,": <START LIST OF CLOSEST CONTACTS INTERNALLY TO ",target.Address,">")

	//expectedWorkCount += len(contacts.contacts)
	for _, contact := range contacts.contacts { //Send request to nodes own closests contacts for their closests contacts
	//	log.Println(kademlia.routingTable.me.Address, ": -->", contact.Address)
		expectedWorkCount++
		log.Println(kademlia.routingTable.me.Address,":+++++++ INCREASED EXPECTED WORK COUNT TO:",expectedWorkCount,"\n")
		//go func(goContact Contact) {
			//TODO: Check if adding go below messes it up????
			kademlia.network.SendFindContactMessage(target, &contact, worker.id,[]Contact{},false) //Create and send lookupcontact request to contact
		//}(contact)
	}
	//log.Println(kademlia.routingTable.me.Address,": </END LIST OF CLOSEST CONTACTS INTERNALLY TO ",target.Address,">")

	//log.Println(kademlia.routingTable.me.Address, ":FINISHED LIST OF CLOSEST CONTACTS INTERNALLY")


lookForRepliesChannel:
		for {
			select {
				case reply := <- worker.workRequest: //Idle wait for replies to requests
					workRecievedCount++ //increment received work counter when received reply
					log.Println(kademlia.routingTable.me.Address,":+-+-+-+ INCREASED RECEIVED WORK COUNT TO:",workRecievedCount,"\n")

					//log.Println(kademlia.routingTable.me.Address,":--------------------- Work received count:",workRecievedCount)

					replyContactsProto := reply.message.GetMsg_2().Contacts //Get the Contact Proto message from the reply
					//Convert protbuf contact to kademlia contact
					replyContacts := []Contact{}
					for _, replyContact := range replyContactsProto{
						if replyContact.KademliaId != kademlia.network.routingTable.me.ID.String(){
							replyKademlia := NewKademliaID(replyContact.KademliaId) //Create kademliaId from reply
							replyContact := NewContact(replyKademlia,replyContact.Address) //Create contact from reply
							replyContact.CalcDistance(target.ID) //Calculate distance
							replyContacts = append(replyContacts, replyContact) //Append reply contacts into a list
						}else{
							log.Println(kademlia.routingTable.me.Address,":Filtered from sending [Find Contact] message to self")
						}
					}

					newContacts := contacts.AppendClosestContacts(replyContacts,3) //Add replied contacts to the contact candidate list

				/**	log.Println(kademlia.routingTable.me.Address,":CURRENT CONTACTS:")
					for _, elem := range contacts.contacts{
						log.Println(kademlia.routingTable.me.Address,":",elem.Address)
					}
					log.Println(kademlia.routingTable.me.Address,":POSSIBLE NEW CONTACTS:")
					for _, elem := range replyContacts{
						log.Println(kademlia.routingTable.me.Address,":",elem.Address)
					}
					log.Println("NEW LENGHT: ",len(contacts.contacts))
					log.Println("NEW CONTACTS")
					for _,elem := range newContacts{
						log.Println(kademlia.routingTable.me.Address,":",elem.Address)
					}
					log.Println("ALL CONTACTS")
					for _,elem := range contacts.contacts{
						log.Println(kademlia.routingTable.me.Address,":",elem.Address)
					} */

					for _,contact := range newContacts{
						expectedWorkCount++
						log.Println(kademlia.routingTable.me.Address,":+++++++ INCREASED EXPECTED WORK COUNT TO:",expectedWorkCount,"\n")

						//TODO: SET GO ROUTINE HERE AGAIN
						kademlia.network.SendFindContactMessage(target,&contact,worker.id,[]Contact{},false) //Fire off a new find contact request for each contact
					}

					//log.Println(kademlia.routingTable.me.Address,":---------------------  Work expected count:",expectedWorkCount)

					if workRecievedCount < expectedWorkCount{
					//	log.Println("")
						log.Println(kademlia.routingTable.me.Address,":ADDED worker back to queue, ID: ",worker.id)
						log.Println(kademlia.routingTable.me.Address,": WORK RECEIVED COUNT:",workRecievedCount)
						log.Println(kademlia.routingTable.me.Address,": EXPECTED WORK COUNT:",expectedWorkCount,"\n")
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

func (kademlia *Kademlia) LookupData(targetHash string) {
	targetId := NewKademliaID(targetHash)
	workRecievedCount, expectedWorkCount := 0, 0 //Keep track of how many requests and replies have been sent
	worker := kademlia.NewWorker() //Identity of the process running
	defer close(worker.workRequest) //Make sure worker is taken out of worker queue when done
	contactCandidates := ContactCandidates{kademlia.routingTable.FindClosestContacts(NewKademliaID(targetHash),3) } //Retrieve nodes own closest contacts

	log.Println(kademlia.routingTable.me.Address,": Started LOOKUP_DATA ", worker.id, " for hash ", targetHash)

	go func(contacts []Contact, count *int) { //Send requests concurrently
		for _, contact := range contacts {
			kademlia.network.SendFindDataMessage(*targetId, contact, worker.id, false, false, nil, nil)
			*count++
		}
	}(contactCandidates.contacts, &workRecievedCount)
	
	timer := time.NewTimer(time.Second * 10) //Timer for timeout, can be reset when replies are received
	LookupDataLoop:
	for {
		select {
			case reply := <- worker.workRequest: //Idle wait for replies to requests
				timer.Reset(time.Second * 10) //Reset because there is further hope
				workRecievedCount++
				
				//See if the reply contains the address to the where the data is located
				if reply.message.GetMsg_3().GetFoundFile() {
					// TODO Store data
					// TODO Create and send store request to next closest node(s?)
					log.Println(kademlia.routingTable.me.Address,": LOOKUP_DATA ", worker.id, " downloaded and published data")
					break LookupDataLoop //Process finished with what it intended to do
					
				//Otherwise, see what can be done with the replies' closest contacts instead
				} else {
					//Extract contacts from reply and filter out self
					replyContacts := ConvertProtobufContacts(reply.message.GetMsg_3().GetContacts(), kademlia.routingTable.me)
					
					//Attempt to add the contacts from the reply to candidates and save the ones added
					newContactCandidates := contactCandidates.AppendNonDuplicates(replyContacts)
					
					//See if there are any new candidates to send more requests to
					if len(newContactCandidates) > 0 {
						go func(contacts []Contact, count *int) { //Send requests concurrently
							for _, contact := range contacts {
								kademlia.network.SendFindDataMessage(*targetId, contact, worker.id, false, false, nil, nil)
								*count++
							}
						}(newContactCandidates, &workRecievedCount)

					//See if it was the last reply, in that case the process failed to find the data
					} else if expectedWorkCount == workRecievedCount {
						log.Println(kademlia.routingTable.me.Address,": LOOKUP_DATA ", worker.id, " couldn't find the data it was looking for")
						break LookupDataLoop //Process didn't find the data it was looking for
					}
					
					//There's more replies for requests out there to handle, better get back to work
					kademlia.network.WorkerQueue <- worker
				}

			//Idle wait for a timeout, timeout resets when a reply arrives
			case <- timer.C:
				log.Println(kademlia.routingTable.me.Address,": Lookup data process ", worker.id, " timed out from lack of replies")
				break LookupDataLoop //Process timed out from lack of replies
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
