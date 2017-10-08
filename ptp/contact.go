package ptp

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
* Append 2 lists together, sort them, keep only the "maxSize" closest contacts
  return the new contacts that "made the cut" and didn't already exist
 */
func (candidates *ContactCandidates) AppendClosestContacts(contacts []Contact, targetId KademliaID, me Contact) []Contact {
	//Keep track of added contacts
	addedCandidates := []Contact{}

	//Make sure all contacts have distances and both arrays are sorted
	newCandidates := ContactCandidates{ contacts }
	for i,_ := range candidates.contacts { candidates.contacts[i].CalcDistance(&targetId) }
	for i,_ := range newCandidates.contacts { newCandidates.contacts[i].CalcDistance(&targetId) }
	candidates.Sort()
	newCandidates.Sort()

	//Go through the new candidates
	count := 0
	NewCandidates:
	for _,newCandidate := range newCandidates.contacts {

		//See if it's a duplicate or own id, skip it if it is
		for _,candidate := range candidates.contacts {
			if newCandidate.ID.Equals(candidate.ID) || newCandidate.ID.Equals(me.ID) {
				continue NewCandidates
			}
		}

		candidates.Append([]Contact{ newCandidate })
		addedCandidates = append(addedCandidates, newCandidate)
		count++

		if count > 3 {
			break
		}
	}

	//Return the added candidates
	return addedCandidates
}

//Append contacts without duplicates
//Returns the added contacts
func (candidates *ContactCandidates) AppendNonDuplicates(contacts []Contact) []Contact{
	nonDupNewContacts := []Contact{}

	for _, contact := range contacts{
		exist := false
		for _, existingElem := range candidates.contacts{
			if contact.ID.String() == existingElem.ID.String() {
				exist = true
				break
			}
		}
		if !exist {
			nonDupNewContacts = append(nonDupNewContacts, contact)
		}
	}

	candidates.Append(nonDupNewContacts)

	return nonDupNewContacts
}

func (candidates *ContactCandidates) FilterClosest(targetId *KademliaID, count int) {
	if len(candidates.contacts) > count {
		for i,_ := range candidates.contacts {
			candidates.contacts[i].CalcDistance(targetId)
		}
		candidates.Sort()
		candidates.contacts = candidates.contacts[:count]
	}
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

func (candidates *ContactCandidates) Remove(contact Contact) {
	filteredCandidates := ContactCandidates{ []Contact{} }
	for _,contactCandidate := range candidates.contacts {
		if !contactCandidate.ID.Equals(contact.ID) {
			filteredCandidates.Append([]Contact{ contactCandidate })
		}
	}
	*candidates = filteredCandidates
}
