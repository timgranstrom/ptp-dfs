package ptp

const bucketSize = 20 //contains a k number of nodes.

/**
*Create a Routing table struct with a Contact (me)
* and an array (buckets) of IDLength*8 length.
 */
type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket //IDLength*8 because IDLength represent byte, hence why we multiply by 8 to get a total of 160 bits.
}
/**
* Pass a contact (me) as argument
* Creates a new routing table and populates the buckets array with new buckets.
* The routing table is assigned the contact.
 */
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}
/**
* Adds a Contact (contact) to a routing table
 */
func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID) //Get bucket index for contact in routingTable.
	bucket := routingTable.buckets[bucketIndex] //Get the bucket.
	bucket.AddContact(contact) //Add the contact to the bucket.
}
/**
* Find closest contacts from kademliaID (target) in routingTable
 */
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target) //Get bucket index of the KademliaID (target)
	bucket := routingTable.buckets[bucketIndex] //Retrieve bucket

	//Get all contacts in the bucket and calc the distance between them and target.
	//Then append these contacts into candidates.
	candidates.Append(bucket.GetContactAndCalcDistance(target))

	//Loop through all buckets or until the nr of candidates have been reached.
	//Start from the targets bucket.
	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i] //Get previous bucket
			candidates.Append(bucket.GetContactAndCalcDistance(target)) //Get all contacts in bucket with distance to target
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i] //Get next bucket
			candidates.Append(bucket.GetContactAndCalcDistance(target)) //Get all contacts in bucket with distance to target
		}
	}

	candidates.Sort() //Sort candidates in increasing distance

	if count > candidates.Len() {
		count = candidates.Len() //Set the count to match the amount of candidates
	}

	return candidates.GetContacts(count) //Returns a list of contacts from candidates.
}

/**
* Get the a bucket index from a routing table for a specific kademliaId.
 */
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID) //get distance array from given ID, to my ID.
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j //Bitshift to find the index of the bucket and return it.
			}
		}
	}

	return IDLength*8 - 1
}
