package ptp

import (
	"testing"
	"time"
	"io/ioutil"
	"log"
	"encoding/hex"
	"container/list"
	"fmt"
	"go/ast"
)

func CreateAndRunNodes(amount int) list.List {
	nodeList := list.New()
	for i := 0; i < amount; i++ {
		nodeList.PushBack(Kademlia{})
		if i == 0 {
			nodeList.Back().Value = NewKademlia(fmt.Sprintf(":%04d", i), nil)
		} else {
			bootstrap := nodeList.Back().Prev().Value.(Kademlia).routingTable.me
			nodeList.Back().Value = NewKademlia(fmt.Sprintf(":%04d", i), &bootstrap)
		}
		go nodeList.Back().Value.(Kademlia).Run()
	}
	return *nodeList
}

func TestRunKademliaInstances(t *testing.T){
	//Create nodes
	node1 := NewKademlia(":8001",nil) //Original node
	//node1ContactCopy := node1.routingTable.me
	//node1ContactCopy.Address = "127.0.0.1"+node1ContactCopy.Address

	node2 := NewKademlia(":8002",&node1.routingTable.me) //boostrap on node1
	//node2ContactCopy := node2.routingTable.me
	//node2ContactCopy.Address = "127.0.0.1"+node2ContactCopy.Address

	node3 := NewKademlia(":8003", &node2.routingTable.me) //boostrap on node2

	node4 := NewKademlia(":8004", &node3.routingTable.me) //boostrap on node2

	//Run nodes
	println("1.----------------------------\n\n")
	//go node1.network.Listen()
	go node1.Run()
	time.Sleep(time.Second)
	println("2.----------------------------\n\n")

	//go node2.network.Listen()
	go node2.Run()
	time.Sleep(time.Second)
	println("3.------------------------------\n\n")
	//go node2.LookupContact(&node1.routingTable.me)
	//println("4.------------------------------\n\n")
	//go node1.LookupContact(&node2.routingTable.me)
	time.Sleep(time.Second)

	//go node1.network.Listen()
	go node3.Run()
	time.Sleep(time.Second)
	println("4.----------------------------\n\n")
	go node4.Run()
	time.Sleep(time.Second)
	println("node2.LookcupContact node 1 ----------------------------\n\n")
	//Try and find node1 through the network
	go node2.LookupContact(&node1.routingTable.me)
	//time.Sleep(time.Second*1)
	println("node3.LookcupContact node 2 ----------------------------\n\n")

	node3.LookupContact(&node2.routingTable.me)
	//time.Sleep(time.Second*10)

	//time.Sleep(time.Second)*/
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
	node1 := NewKademlia(":8001", nil) //Original node
	node2 := NewKademlia(":8002", &node1.routingTable.me) //Original node
	node3 := NewKademlia(":8003", &node2.routingTable.me) //Original node
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
		t.FailNow()
	}
	time.Sleep(time.Second)

	go node3.Run()
	time.Sleep(time.Second)
	if node3.LookupData(hashKey) {
		log.Println("DATA FOUND, TEST COMPLETE")
	} else {
		log.Println("DATA NOT FOUND, TEST FAILED")
		t.Fail()
	}
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
