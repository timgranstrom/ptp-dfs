package ptp

import (
	"container/list"
)

/*
*Create a bucket which has a list.
 */
type bucket struct {
	list *list.List
}

/**
* Create a new bucket and give it a new list.
 */
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

/**
* Add a contact to a bucket.
 */
func (bucket *bucket) AddContact(contact Contact) *Contact {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() { //Iterate through list of contacts
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e //if contact already exist in bucket, assign that contact (e) to element.
		}
	}
	//if contact doesn't exist in bucket and bucket is not full, add contact to front of the list.
	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {
			return bucket.list.Back().Value.(*Contact)
		}
	} else {
		bucket.list.MoveToFront(element) //If contact already exist in bucket, only move it to the front.
	}
	return nil
}
/*
* Get contacts and calculate distance in bucket from target.
 */
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact) //Get element from list and typecast into a Contact
		contact.CalcDistance(target) //Get distance between current contact and target
		contacts = append(contacts, contact) //Append current contact into contacts.
	}

	return contacts
}
/**
* Get length of the bucket.
 */
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
