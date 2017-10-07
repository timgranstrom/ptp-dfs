package ptp

import (
	"testing"
	"time"
	"io/ioutil"
	"log"
	"encoding/hex"
)

func TestDaemonService_CatFromNetwork(t *testing.T)  {

	//Create nodes
	ds := NewDaemonService() //daemon node
	node1 := NewKademlia(":8001", &ds.kademliaNode.routingTable.me) //node 1
	node2 := NewKademlia(":8002", &node1.routingTable.me) //node 2
	node3 := NewKademlia(":8003", &node2.routingTable.me) //node 3
	node4 := NewKademlia(":8004", &node3.routingTable.me) //node 4
	node5 := NewKademlia(":8005", &node4.routingTable.me) //node 5

	time.Sleep(time.Second)
	go ds.kademliaNode.Run()
	time.Sleep(time.Second)
	go node1.Run()
	time.Sleep(time.Second)
	go node2.Run()
	time.Sleep(time.Second)
	go node3.Run()
	time.Sleep(time.Second)
	go node4.Run()
	time.Sleep(time.Second)
	go node5.Run()
	time.Sleep(time.Second)

	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	b, _ := ioutil.ReadFile(path) // Take out the content of the file in byte
	hashKey := node5.Store(fileName,b)
	time.Sleep(time.Second)
	content := ds.Cat(hashKey)
	if content != "" {
		log.Println("Content: ",content)
		log.Println("CAT DATA FOUND, TEST COMPLETE")
	} else {
		log.Println("CAT DATA NOT FOUND, TEST FAILED")
		t.Fail()
	}
}

func TestDaemonService_CatLocal(t *testing.T)  {

	//Create nodes
	ds := NewDaemonService() //daemon node
	time.Sleep(time.Second)
	go ds.kademliaNode.Run()
	time.Sleep(time.Second)

	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	b, _ := ioutil.ReadFile(path) // Take out the content of the file in byte
	hashKey := ds.kademliaNode.Store(fileName,b)
	time.Sleep(time.Second)
	content := ds.Cat(hashKey)
	if content != "" {
		log.Println("Content: ",content)
		log.Println("CAT DATA FOUND LOCALLY, TEST COMPLETE")
	} else {
		log.Println("CAT DATA NOT FOUND, TEST FAILED")
		t.Fail()
	}
}

func TestDaemonService_Pin(t *testing.T)  {

	//Create nodes
	ds := NewDaemonService() //daemon node
	time.Sleep(time.Second)
	go ds.kademliaNode.Run()
	time.Sleep(time.Second)

	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	b, _ := ioutil.ReadFile(path) // Take out the content of the file in byte
	hashKey := ds.kademliaNode.Store(fileName,b)
	time.Sleep(time.Second)
	hexKey,_ := hex.DecodeString(hashKey)
	ds.Pin(hashKey,true)
	isPinned, _ := ds.kademliaNode.network.store.RetrieveIsPinned(hexKey)
	if isPinned {
		log.Println("Success: Data is pinned.")
	} else {
		log.Println("Fail: Data is not pinned.")
		t.Fail()
	}
}

func TestDaemonService_Store(t *testing.T)  {

	//Create nodes
	ds := NewDaemonService() //daemon node
	node1 := NewKademlia(":8001", &ds.kademliaNode.routingTable.me) //node 1
	node2 := NewKademlia(":8002", &node1.routingTable.me) //node 2
	node3 := NewKademlia(":8003", &node2.routingTable.me) //node 3
	node4 := NewKademlia(":8004", &node3.routingTable.me) //node 4
	node5 := NewKademlia(":8005", &node4.routingTable.me) //node 5

	time.Sleep(time.Second)
	go ds.kademliaNode.Run()
	time.Sleep(time.Second)
	go node1.Run()
	time.Sleep(time.Second)
	go node2.Run()
	time.Sleep(time.Second)
	go node3.Run()
	time.Sleep(time.Second)
	go node4.Run()
	time.Sleep(time.Second)
	go node5.Run()
	time.Sleep(time.Second)

	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	b, _ := ioutil.ReadFile(path) // Take out the content of the file in byte
	hashKey := ds.kademliaNode.Store(fileName,b)
	time.Sleep(time.Second)
	hexKey,_ := hex.DecodeString(hashKey)
	_,isFound := ds.kademliaNode.network.store.RetrieveData(hexKey)
	if isFound {
		log.Println("Success: Data is stored.")
	} else {
		log.Println("Fail: Data is not stored.")
		t.Fail()
	}
}
