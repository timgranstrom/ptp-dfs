package ptp

import (
	"encoding/hex"
	"math/rand"
)

const IDLength = 20

type KademliaID [IDLength]byte

/*
*
 */
func NewKademliaID(data string) *KademliaID {
	//Returns the bytes represented by the hexadecimal string data.
	decoded, _ := hex.DecodeString(data) //decoded = the byte, _ = ignore the error returned if any.

	newKademliaID := KademliaID{} //Create a new empty KademliaID
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i] //Add bytes to KademliaID.
	}

	return &newKademliaID //Return the Address of the kademliaID.
}
/**
* Generate a randomized KademliaID.
 */
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256)) //Assign random bits to the KademliaID
	}
	return &newKademliaID //Return address of the kademliaID.
}
//Find the lesser KademliaID.
//returns true if kademliaID is less than OtherKademliaId.
//Otherwise return false.
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}
/*
* Check if two KademliaID are equal to each other.
 */
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

/**
* Calculate distance between 2 nodes and return the path/distance between them in bits.
 */
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}
/**
* Returns the string representation of the KademliaId.
 */
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
