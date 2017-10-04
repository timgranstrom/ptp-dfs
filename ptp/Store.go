package ptp

import (
	"time"
	"sync"
	"crypto/sha1"
)

/*
Store structure that has a key/value map for data
 */
type Store struct{
	mutex *sync.Mutex
	data map[string][]byte
	expirationTime map[string]time.Time
	republishTime map[string]time.Time
	lifeTimeDuration map[string]time.Duration
	pin map[string]bool
}

// Init initializes the Store
func MakeStore() *Store {
	store := Store{}
	store.data = make(map[string][]byte)
	store.mutex = &sync.Mutex{}
	store.expirationTime = make(map[string]time.Time)
	store.republishTime = make(map[string]time.Time)
	store.lifeTimeDuration = make(map[string]time.Duration)
	store.pin = make(map[string]bool)
	return &store
}

/**
 Retrieve data from the store using a key
 */
func (store *Store) RetrieveData(key []byte) (data []byte, isFound bool){
	store.mutex.Lock()
	defer store.mutex.Unlock()
	data, isFound = store.data[string(key)] //Get data that matches the key from the store
	return data,isFound
}

/**
Retrieve the expiration time for a key from the store
 */
func (store *Store) RetrieveExpirationTime(key []byte) (time time.Time, isFound bool){
	store.mutex.Lock()
	defer store.mutex.Unlock()
	time, isFound = store.expirationTime[string(key)] //Get data that matches the key from the store
	return time,isFound
}

/**
Retrieve the expiration time for a key from the store
 */
func (store *Store) RetrieveRepublishTime(key []byte) (time time.Time, isFound bool){
	store.mutex.Lock()
	defer store.mutex.Unlock()
	time, isFound = store.republishTime[string(key)] //Get data that matches the key from the store
	return time,isFound
}

/**
Retrieve the expiration time for a key from the store
 */
func (store *Store) RetrieveLifeTimeDuration(key []byte) (time time.Duration, isFound bool){
	store.mutex.Lock()
	defer store.mutex.Unlock()
	time, isFound = store.lifeTimeDuration[string(key)] //Get data that matches the key from the store
	return time,isFound
}

/**
Retrieve info if the key/data is pinned or not.
 */
func (store *Store) RetrieveIsPinned(key []byte) (isPinned bool, isFound bool){
	store.mutex.Lock()
	defer store.mutex.Unlock()
	isPinned, isFound = store.pin[string(key)] //Get data that matches the key from the store
	return isPinned,isFound
}

/**
Retrieve info if the key/data is pinned or not.
 */
func (store *Store) SetPin(key []byte,isPinned bool) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.pin[string(key)] = isPinned //Mark the key as pinned or not pinned
	return nil
}

// GetKey returns the key for data
func (store *Store) GetKey(fileName string) []byte {
	sha := sha1.Sum([]byte(fileName))
	return sha[:]
}

/**
Store data and expiration time for a specific key
 */
func (store *Store) StoreData(key []byte, data []byte, lifetimeDuration time.Duration, republishTime time.Time,pinned bool) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.data[string(key)] = data
	store.expirationTime[string(key)] = time.Now().Add(lifetimeDuration)
	store.republishTime[string(key)] = republishTime
	store.pin[string(key)] = pinned
	return nil
}

/**
 Delete data and expiration time for a specific key
 */
func (store *Store) Delete(key []byte) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()
	isPinned, isFound := store.pin[string(key)] //Get data that matches the key from the store
	//Remove all data if the key is not pinned
	if isFound && !isPinned{
		delete(store.data, string(key))           //delete data for specific key
		delete(store.expirationTime, string(key)) //delete expiration time for specific key
		delete(store.lifeTimeDuration, string(key)) //delete lifetime time for specific key
		delete(store.republishTime, string(key)) //delete republish time for specific key
		delete(store.pin, string(key)) //delete expiration time for specific key
	}
	return nil
}