package ptp

import (
	"testing"
	"time"
)

func TestRunKademliaInstances(t *testing.T){
	//Create nodes
	node1 := NewKademlia(":8001",nil) //Original node
	node2 := NewKademlia(":8002",&node1.routingTable.me) //boostrap on node1
	node3 := NewKademlia(":8003", &node2.routingTable.me) //boostrap on node2
	//Run nodes
	node1.Run()
	node2.Run()
	node3.Run()
	time.Sleep(time.Second)

	//Try and find node1 through the network
	go node3.LookupContact(&node1.routingTable.me)

	time.Sleep(time.Second)
}
