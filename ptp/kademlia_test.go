package ptp

import (
	"testing"
	"time"
	"io/ioutil"
	"log"
	"encoding/hex"
	"container/list"
	"fmt"
)

func CreateAndRunNodes(amount int, originContact *Contact) list.List {
	nodeList := list.New()
	for i := 0; i < amount; i++ {
		nodeList.PushBack(Kademlia{})
		if i == 0 {
			nodeList.Back().Value = NewKademlia(fmt.Sprintf(":9%03d", i), originContact)
		} else {
			bootstrap := nodeList.Back().Prev().Value.(*Kademlia).routingTable.me
			nodeList.Back().Value = NewKademlia(fmt.Sprintf(":9%03d", i), &bootstrap)
		}
		node := nodeList.Back().Value.(*Kademlia)
		go node.Run()
		time.Sleep(time.Second * 3)
	}
	return *nodeList
}

func TestKademliaStartAndRunNodes(t *testing.T){
	CreateAndRunNodes(100,nil)
}

func TestKademliaWithDaemon(t *testing.T){
	daemonContact := NewContact(NewKademliaID("210fc7bb818639ac48a4c6afa2f1581a8b9525e2"),":8000")
	CreateAndRunNodes(10,&daemonContact)
	time.Sleep(time.Second*1000)
}

func TestStoreKademlia(t *testing.T) {
	//Create nodes
	node1 := NewKademlia(":8001", nil) //Original node
	node2 := NewKademlia(":8002", &node1.routingTable.me) //Original node
	time.Sleep(time.Second)
	go node1.Run()
	time.Sleep(time.Second)
	go node2.Run()
	time.Sleep(time.Second)

	ds := NewDaemonService()
	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	b, _ := ioutil.ReadFile(path) // Take out the content of the file in byte
	hashKey := node1.Store(fileName,b)
	time.Sleep(time.Second)
	data,_ := hex.DecodeString(hashKey)
	data,isFound := node2.network.store.RetrieveData(data)
	if isFound{
		log.Println("DATA CONTENT: \n",string(data))
	} else{
		log.Fatal("Did not find expected file")
		t.Fail()
	}
	time.Sleep(time.Second)
}

func TestPingKademlia(t *testing.T) {
	//Create nodes
	node1 := NewKademlia(":8001", nil) //Original node
	node2 := NewKademlia(":8002", &node1.routingTable.me) //Original node
	time.Sleep(time.Second)
	go node1.Run()
	time.Sleep(time.Second)
	go node2.Run()
	time.Sleep(time.Second)

	if node1.PingContact(Ping{ make(chan bool), &node2.routingTable.me}) {
		log.Println("Successfully pinged the other contact")
	} else {
		log.Println("Failed to ping contact, test failed")
		t.Fail()
	}
}

func TestLookupDataKademlia(t *testing.T)  {
	//Create nodes
	nodeList := CreateAndRunNodes(100,nil)

	//Retrieve data
	ds := NewDaemonService()
	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	data, _ := ioutil.ReadFile(path) // Take out the content of the file in byte

	//Store data in the first made node
	frontNode := nodeList.Front().Value.(*Kademlia)
	hashKey := frontNode.Store(fileName, data)
	time.Sleep(time.Second)
	println()

	//Look it up from the last made node
	backNode := nodeList.Back().Value.(*Kademlia)
	foundData, isFound := backNode.LookupData(hashKey)

	//See if the last node could find the data the first node stored
	if isFound {
		println(string(foundData))
	} else {
		t.Fail()
	}

	time.Sleep(time.Second * 30)
}

func TestPurgeDataKademlia(t *testing.T) {
	//Create nodes
	node1 := NewKademlia(":8001", nil) //Original node
	node2 := NewKademlia(":8002", &node1.routingTable.me) //Original node
	time.Sleep(time.Second)
	go node1.Run()
	time.Sleep(time.Second)
	go node2.Run()
	time.Sleep(time.Second)

	ds := NewDaemonService()
	fileName, path := ds.ParseFilePathCommand("../main/file.txt")
	b, _ := ioutil.ReadFile(path) // Take out the content of the file in byte
	node1.Store(fileName,b) //Try to store, check if they republish automatically (REMEMBER TO CHANGE REPUBLISH TIME FOR TESTING
	time.Sleep(time.Second*60)
}
