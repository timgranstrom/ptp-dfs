package ptp

import (
	"time"
	"sync"
	"crypto/sha1"
	"fmt"
)

const(
	ExpirationTimeOriginal = time.Second*24+time.Second*6 //Should be 24 hours and 6 minutes (just a little while after republish) ||| is 30 seconds for testing purposes
	ExpirationTimeReplicate = time.Second*10+time.Second*6 //Should be 1 hour and 6 minutes (just a little while after republish) ||| is 16 seconds for testing purposes
	RepublishTimeOriginal = time.Second*24 //Should be 24 hours ||| is 24 seconds for testing purposes
	RepublishTimeReplicate = time.Second*10 //Should be 1 hour ||| is 10 seconds for testing purposes
)

/*
Store structure that has a key/value map for data
 */
type Store struct{
	mutex *sync.Mutex
	storeObjects map[string]StoreObject //Stored objects
	//data map[string][]byte
	//expirationTime map[string]time.Time
	//republishTime map[string]time.Time
	//lifeTimeDuration map[string]time.Duration
	//pin map[string]bool
}

type StoreError struct {
	s string
}
func (e *StoreError) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
func NewStoreError(text string) error {
	return &StoreError{text}
}

// Init initializes the Store
func MakeStore() *Store {
	store := Store{}
	store.storeObjects = make(map[string]StoreObject)
	store.mutex = &sync.Mutex{}
	//store.data = make(map[string][]byte)
	//store.expirationTime = make(map[string]time.Time)
	//store.republishTime = make(map[string]time.Time)
	//store.lifeTimeDuration = make(map[string]time.Duration)
	//store.pin = make(map[string]bool)
	return &store
}

type StoreObject struct{
	expirationTime time.Time
	republishTime time.Time
	pinned bool
	data []byte
	isOriginal bool
}

func NewStoreObject(data[]byte,isOriginal bool) *StoreObject{
	if isOriginal{
		storeObject := &StoreObject{expirationTime:time.Now().Add(ExpirationTimeOriginal),
		republishTime:time.Now().Add(RepublishTimeOriginal),
		pinned: false,
		data:data,
		isOriginal:isOriginal}
		return storeObject
	}
	storeObject := &StoreObject{expirationTime:time.Now().Add(ExpirationTimeReplicate),
		republishTime:time.Now().Add(RepublishTimeReplicate),
		pinned: false,
		data:data,
		isOriginal:isOriginal}

	return storeObject
}

func (storeObject *StoreObject) ResetExpirationTime(){
	if storeObject.isOriginal{
		fmt.Println("RESET Expiration Time for original")

		storeObject.expirationTime = time.Now().Add(ExpirationTimeOriginal)
	} else{
		fmt.Println("RESET expiration Time for replicate")
		storeObject.expirationTime = time.Now().Add(ExpirationTimeReplicate)
	}
}

func (storeObject *StoreObject) ResetRepublishTime(){
	if storeObject.isOriginal{
		fmt.Println("RESET Republish Time for original")
		storeObject.republishTime = time.Now().Add(RepublishTimeOriginal)
	} else{
		fmt.Println("RESET Republish Time for replicate")
		storeObject.republishTime = time.Now().Add(RepublishTimeReplicate)
	}
}

/**
 Retrieve a StoreObject from the store using a key
 */
func (store *Store) RetrieveStoreObject(key []byte) (storeObject StoreObject, isFound bool){
	store.mutex.Lock()
	defer store.mutex.Unlock()
	storeObject, isFound = store.storeObjects[string(key)] //Get store object that matches a key
	return storeObject,isFound
}

/**
 Retrieve data from the store using a key
 */
func (store *Store) RetrieveData(key []byte) (data []byte, isFound bool){
	//store.mutex.Lock()
	//defer store.mutex.Unlock()
	storeObject,foundStoreObj := store.RetrieveStoreObject(key)
	if foundStoreObj{
		return storeObject.data,foundStoreObj
	}
	return nil,foundStoreObj
}

/**
Retrieve the expiration time for a key from the store
 */
func (store *Store) RetrieveExpirationTime(key []byte) (time time.Time, isFound bool){
	//store.mutex.Lock()
	//defer store.mutex.Unlock()
	storeObject,foundStoreObj := store.RetrieveStoreObject(key)
	return storeObject.expirationTime,foundStoreObj
}

/**
Retrieve the expiration time for a key from the store
 */
func (store *Store) RetrieveRepublishTime(key []byte) (time time.Time, isFound bool){
	//store.mutex.Lock()
	//defer store.mutex.Unlock()
	storeObject,foundStoreObj := store.RetrieveStoreObject(key)
	return storeObject.republishTime,foundStoreObj
}

/**
Retrieve info if the key/data is pinned or not.
 */
func (store *Store) RetrieveIsPinned(key []byte) (isPinned bool, isFound bool){
	//store.mutex.Lock()
	//defer store.mutex.Unlock()
	storeObject,foundStoreObj := store.RetrieveStoreObject(key)
	return storeObject.pinned,foundStoreObj
}

/**
Retrieve info if the key/data is pinned or not.
 */
func (store *Store) SetPin(key []byte,isPinned bool) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()
	//store.pin[string(key)] = isPinned //Mark the key as pinned or not pinned
	storeObject, isFound := store.storeObjects[string(key)] //Get store object that matches a key
	if isFound{
		storeObject.pinned = isPinned
		store.storeObjects[string(key)] = storeObject
		return nil
	}
	return NewStoreError("Could not find store object")
}

// GetKey returns the key for data
func (store *Store) GetKey(fileName string) []byte {
	sha := sha1.Sum([]byte(fileName))
	return sha[:]
}

/**
Store data and expiration time for a specific key
 */
func (store *Store) StoreData(key []byte, data []byte, isOriginal bool) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()
	storeObject,alreadyExist := store.storeObjects[string(key)]

	if !alreadyExist{
		//fmt.Println("STORED NEW OBJECT")
		storeObject = *NewStoreObject(data,isOriginal) //If it doesn't already exist, just make it
	} else{ //If it already exist, reset it depending if it is has the original data or not
		storeObject.ResetRepublishTime()
		storeObject.ResetExpirationTime()
	}
	//store.data[string(key)] = data
	//store.expirationTime[string(key)] = time.Now().Add(lifetimeDuration)
	//store.republishTime[string(key)] = republishTime
	//store.pin[string(key)] = pinned

	store.storeObjects[string(key)] = storeObject
	return nil
}

/**
 Delete data and expiration time for a specific key
 */
func (store *Store) Delete(key []byte) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()
	storeObject,isFound := store.storeObjects[string(key)] //Get data that matches the key from the store
	//Remove all data if the key is not pinned
	if isFound && !storeObject.pinned{
		delete(store.storeObjects,string(key)) //Delete store object for specific key
		fmt.Println("DELETED AN OBJECT")
	}
	if !isFound{
		return NewStoreError("Could not find store object")
	}
	return nil
}