package ptp

import (
	"testing"
	"time"
	"io/ioutil"
	"log"
)

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

	storeLocation := "main/file.txt"
	b, _ := ioutil.ReadFile(storeLocation) // Take out the content of the file in byte
	hashKey := node1.Store("file.txt",b)
	time.Sleep(time.Second)
	data,isFound := node2.network.store.RetrieveData(hashKey)
	if isFound{
		log.Println("DATA CONTENT: \n",string(data))
	} else{
		log.Fatal("Did not find expected file")
	}
	time.Sleep(time.Second)

	//node1.Store()
}
