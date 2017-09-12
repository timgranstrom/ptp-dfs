package ptp

type Kademlia struct {
	routingTable RoutingTable
	network Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//alpha := 3 //count = alpha = 3 is standard nr.
	//contacts := kademlia.routingTable.FindClosestContacts(target.ID,alpha)
	//for _, element := range contacts{
		//Start new goroutine for each contact
		//Listen for response
		//Each new, add the closest ones.
	//}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
