package ptp

import (
	"testing"
	"time"
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
	//Run nodes
	go node1.Run()
	time.Sleep(time.Second)
	go node2.Run()
	time.Sleep(time.Second)
	go node3.Run()
	time.Sleep(time.Second/5)
	println("---------------------------------------------------\n\n")
	//Try and find node1 through the network
	//go node2.LookupContact(&node1.routingTable.me)

	//time.Sleep(time.Second*5)
}
