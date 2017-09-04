package ptp_dfs

import (
	"fmt"
	"sort"
)

/**
*Create a contact with a KademliaID, Address and distance
 */
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}
/*
* Creates a new contact with id, address and nil distance.
 */
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}
/*
* Set distance between 2 nodes and return the path/distance between them in bits.
 */
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

/*
* Check if contact distance is less than other contact distance
 */
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}
/*
* returns a formatted string representation like: Contact("ID","ADDRESS")
 */
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

/**
* List of contacts.
 */
type ContactCandidates struct {
	contacts []Contact
}

/**
* Append 2 lists together
 */
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}
/**
* Return the {count} first elements of candidates.contacts.
 */
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}
/*
* Sort candidates in increasing distance
 */
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}
/*
* Get the amount of candidates.
 */
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}
/*
* Swap places of candidates of index i and j.
 */
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}
/**
* return true if contact i is less than contact j.
 */
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}
