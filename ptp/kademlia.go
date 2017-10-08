package ptp

import (
	"container/list"
	"log"
	"time"
	"encoding/hex"
	"fmt"
)

type Kademlia struct {
	routingTable  RoutingTable
	network       *Network
	workers       *list.List //List of all current workers
	idCount       int64      //Global id counter for workers
	bootstrapNode *Contact   //boostrap node to get you into the network
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
		//log.Println(kademlia.routingTable.me.Address+": Add boostrap node as contact")
		kademlia.routingTable.AddContact(*bootstrapNode) //Add boostrap node in network
	}
	return kademlia
}

func (kademlia *Kademlia) Run(){
	println()
	log.Println(kademlia.GetMe().Address + ": Initializing node")
	//println("\n---------------")
	//log.Println(kademlia.network.routingTable.me.Address,": NODE STARTED")

	go kademlia.network.Listen() //Start listener on network
	go kademlia.network.Sender() //Start the sender on network
	dispatcher := NewDispatcher(kademlia.network)
	dispatcher.StartDispatcher() //always run

	go kademlia.BootstrapProcess()

	go kademlia.StoreTimer()
	//See if the routingtable wants to ping a contact
	go func() {
		for {
			select {
				case ping := <- kademlia.network.PingQueue:
					go kademlia.PingContact(ping)
			}
		}
	}()
}

func (kademlia *Kademlia) BootstrapProcess(){
	for {
		if !kademlia.network.listenerActive {
			time.Sleep(time.Second / 10)
		} else{
			break
		}
	}
	log.Println(kademlia.network.routingTable.me.Address + ": Bootstrapping")
	kademlia.LookupContact(&kademlia.routingTable.me) //Find closest nodes to me in network
	log.Println(kademlia.network.routingTable.me.Address + ": Finished bootstrapping")
}

func (kademlia *Kademlia) LookupContact(target *Contact) ContactCandidates {
	workReceivedCount, expectedWorkCount := 0, 0
	worker := kademlia.NewWorker()
	defer worker.SetInactive()
	kademlia.network.WorkerQueue <- worker

	log.Println(kademlia.GetMe().Address + ": LOOKUP_CONTACT", worker.id, "started for", target.ID)

	contactCandidates := ContactCandidates{kademlia.routingTable.FindClosestContacts(target.ID, 3)}

	//Send requests in goroutine so process can instantly wait for replies
	expectedWorkCount += len(contactCandidates.contacts)
	go func(contacts []Contact) {
		for _, contact := range contacts {
			kademlia.network.SendFindContactMessage(target, &contact, worker.id, []Contact{}, false)
			log.Println(kademlia.GetMe().Address + ":  Sending LOOKUP_CONTACT request to  ", fmt.Sprintf("%20s", contact.Address), "for", target.ID)
		}
	}(contactCandidates.contacts)

	timeoutTimer := time.NewTimer(time.Second * 3)

	//Start waiting for replies or timeout
	LookupContactLoop:
	for {
		select {
			case reply := <-worker.workRequest: //Idle wait for replies to requests
				//Reset timeout timer and increment the amount of replies received
				timeoutTimer.Reset(time.Second * 3)
				workReceivedCount++

				//Convert protobuf contacts in reply to kademlia contacts
				replyContacts := ConvertProtobufContacts(reply.message.GetMsg_2().Contacts)

				//Add replied contacts to the contact candidate list, get the newly added ones
				newContacts := contactCandidates.AppendClosestContacts(replyContacts, *target.ID, kademlia.GetMe())

				//Send requests to new closest contacts in goroutine so process can instantly go back to waiting for replies
				expectedWorkCount += len(newContacts)
				go func(contacts []Contact) {
					for _, contact := range contacts {
						kademlia.network.SendFindContactMessage(target, &contact, worker.id, []Contact{}, false)
						log.Println(kademlia.GetMe().Address + ":  Sending LOOKUP_CONTACT request to  ", fmt.Sprintf("%20s", contact.Address), "for", target.ID)
					}
				}(newContacts)

				//If there's more work, go back to waiting, otherwise there's nothing more and closest contacts have been found
				if workReceivedCount < expectedWorkCount {
					kademlia.network.WorkerQueue <- worker
				} else {
					log.Println(kademlia.routingTable.me.Address + ": LOOKUP_CONTACT", worker.id, "found closest contacts, got all expected replies")
					break LookupContactLoop
				}

			//Idle wait for a timeout, timeout resets when a reply arrives
			case <-timeoutTimer.C:
				log.Println(kademlia.routingTable.me.Address + ": LOOKUP_CONTACT", worker.id, "timed out, got", workReceivedCount, "out of", expectedWorkCount, "expected replies")
				break LookupContactLoop
		}
	}
	if contactCandidates.Len() > 3 {
		contactCandidates.Sort()
		contactCandidates.contacts = contactCandidates.contacts[:3]
	}
	return contactCandidates
}

func (kademlia *Kademlia) LookupData(targetHash string) (data []byte, isFound bool) {
	targetId := NewKademliaID(targetHash)
	workRecievedCount, expectedWorkCount := 0, 0 //Keep track of how many requests and replies have been sent
	worker := kademlia.NewWorker() //Identity of the process running
	defer worker.SetInactive()
	kademlia.network.WorkerQueue <- worker //Put self on queue

	contactCandidates := ContactCandidates{kademlia.routingTable.FindClosestContacts(NewKademliaID(targetHash),3) } //Retrieve nodes own closest contacts
	latestNonFileContact := &Contact{} //Contact to send store request to if data is found

	log.Println(kademlia.routingTable.me.Address + ": LOOKUP_DATA", worker.id, "started for hash", targetHash)

	//Send requests in goroutine so process can instantly wait for replies
	expectedWorkCount += len(contactCandidates.contacts)
	go func(contacts []Contact) {
		for _, contact := range contacts {
			log.Println(kademlia.GetMe().Address + ":  Sending LOOKUP_DATA request to  ", fmt.Sprintf("%20s", contact.Address), "for", targetId.String())
			kademlia.network.SendFindDataMessage(targetId, contact, worker.id, false, false, nil, nil)
		}
	}(contactCandidates.contacts)

	//Timer for timeout, can be reset when replies are received
	timeoutTimer := time.NewTimer(time.Second * 3)

	//Start waiting for replies or timeout
	for {
		select {
			case reply := <- worker.workRequest: //Idle wait for replies to requests
				timeoutTimer.Reset(time.Second * 3)
				workRecievedCount++
				
				//See if the reply contains the address to the where the data is located
				if reply.message.GetMsg_3().GetFoundFile() {

					//Extract information to store the data
					data := []byte(reply.message.GetMsg_3().GetFileData())
					key,err := hex.DecodeString(targetHash)
					CheckError(err)
					kademlia.network.store.StoreData(key, data, false) //Store the data

					//See if there is a contact to send a store request to
					if latestNonFileContact != nil {
						kademlia.network.SendStoreMessage(latestNonFileContact, key, data, worker.id, false)
						log.Println(kademlia.GetMe().Address + ":  Sending STORE_DATA request to  ", fmt.Sprintf("%20s", latestNonFileContact.Address))
					}

					log.Println(kademlia.routingTable.me.Address + ": LOOKUP_DATA", worker.id, "found the data")
					return data, true //Process finished with what it intended to do, return true to signify that
					
				//Otherwise, see what can be done with the replies' closest contacts instead
				} else {
					//Save the contact
					latestNonFileContact = &Contact{ NewKademliaID(reply.message.SenderKademliaId), reply.senderAddress, nil }

					//Extract contacts from reply and filter out self
					replyContacts := ConvertProtobufContacts(reply.message.GetMsg_3().GetContacts())
					
					//Attempt to add the contacts from the reply to candidates and save the ones added
					newContactCandidates := contactCandidates.AppendClosestContacts(replyContacts, *targetId, kademlia.GetMe())
					
					//See if there are any new candidates to send more requests to
					expectedWorkCount += len(newContactCandidates)
					go func(contacts []Contact) { //Send requests concurrently
						for _, contact := range contacts {
							log.Println(kademlia.GetMe().Address + ":  Sending LOOKUP_DATA request to  ", fmt.Sprintf("%20s", contact.Address), "for", targetId.String())
							kademlia.network.SendFindDataMessage(targetId, contact, worker.id, false, false, nil, nil)
						}
					}(newContactCandidates)

					//If there's more work, go back to waiting, otherwise there's nothing more and the data couldn't be found
					if workRecievedCount < expectedWorkCount {
						kademlia.network.WorkerQueue <- worker
					} else {
						log.Println(kademlia.routingTable.me.Address + ": LOOKUP_DATA", worker.id, "couldn't find the data it was looking for")
						return nil, false //Process failed to find the data
					}
				}

			//Idle wait for a timeout, timeout resets when a reply arrives
			case <- timeoutTimer.C:
				log.Println(kademlia.routingTable.me.Address + ": LOOKUP_DATA", worker.id, "timed out from lack of replies, recieved", workRecievedCount, "out of", expectedWorkCount, "replies")
				return nil, false //Process failed to find the data due to a timeout
		}
	}
}

/**
Store function to store data from a filename
 */
func (kademlia *Kademlia) Store(fileName string,data []byte) (keyEncoded string) {
	log.Println(kademlia.GetMe().Address + ": STORE_DATA started for", fileName)
	key := kademlia.network.store.GetKey(fileName) //Get the finalized hash result
	keyEncoded = kademlia.StoreInternal(key,data) //Run the internal store function
	log.Println(kademlia.GetMe().Address + ": STORE_DATA finished for", fileName)
	return keyEncoded
}

func (kademlia *Kademlia) StoreInternal(key []byte,data []byte) (keyEncoded string) {
	keyEncoded = hex.EncodeToString(key) //Encode the hash key as a string

	kademlia.network.store.StoreData(key,data,true) //Store data for ourselves as well as the original
	storeKadId := NewKademliaID(keyEncoded) //Make kademlia id out of the key
	storeContact := NewContact(storeKadId,"") //Create contact out of kad id
	contactCandidates := kademlia.LookupContact(&storeContact) //Get the closest contacts to the data
	worker := kademlia.NewWorker() //Just make a worker to get a unique message- and worker id.
	for _,targetContact := range contactCandidates.contacts{
		kademlia.network.SendStoreMessage(&targetContact,key,data,worker.id,false)
		log.Println(kademlia.GetMe().Address + ":  Sending STORE_DATA request to  ", fmt.Sprintf("%20s", targetContact.Address))
	}
	return keyEncoded
}

/**
*Return myself as a contact
 */
func (kademlia *Kademlia) GetMe() Contact{
	return kademlia.routingTable.me
}

func (kademlia *Kademlia) PingContact(ping Ping) bool {
	worker := kademlia.NewWorker()
	defer worker.SetInactive()
	kademlia.network.SendPingMessage(ping.target.Address, worker.id, false)
	select {
		case <- worker.workRequest:
			ping.reply <- true
			return true
		case <- time.NewTimer(time.Second * 3).C:
			ping.reply <- false
			return false
	}
}

func (kademlia *Kademlia) StoreTimer(){

	for{
		var removeStoreKeys []string
		republishStoreObjects := make(map[string]StoreObject)

		kademlia.network.store.mutex.Lock()
		for key,storeObject := range kademlia.network.store.storeObjects{
			if time.Now().After(storeObject.expirationTime){
				removeStoreKeys = append(removeStoreKeys, key) //Flag to remove data for this key
			} else if time.Now().After(storeObject.republishTime){ //If object has not expired, check if it should republish
			fmt.Println("TIME TO REPUBLISH")
				republishStoreObjects[key] = storeObject //save objects to republish and re-store(internally)
			}
		}
		kademlia.network.store.mutex.Unlock()

		for _,key := range removeStoreKeys{
			kademlia.network.store.Delete([]byte(key)) //delete stored data that has expired
		}

		for key,storeObject := range republishStoreObjects{
			nonHexKey := hex.EncodeToString([]byte(key))
			//fmt.Println("TRY TO TRIGGER REPUBLISH FOR ",nonHexKey)
			fmt.Println("TRIGGER REPUBLISH FOR ",nonHexKey)

			kademlia.StoreInternal([]byte(key),storeObject.data)
		}
		time.Sleep(time.Second)
	}

}
