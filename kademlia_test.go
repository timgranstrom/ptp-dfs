package ptp

import (
	"testing"
	"time"
)

func TestRunKademliaInstances(t *testing.T){
	kademliaInstance1 := NewKademlia(":8001")
	kademliaInstance2 := NewKademlia(":8002")
	kademliaInstance3 := NewKademlia(":8003")
	kademliaInstance3.TestAddContact(*kademliaInstance2.TestGetMeContact())
	kademliaInstance2.TestAddContact(*kademliaInstance1.TestGetMeContact())

	kademliaInstance1.Run()
	kademliaInstance2.Run()
	kademliaInstance3.Run()
	time.Sleep(time.Second)
	go kademliaInstance3.TestSendMsg(kademliaInstance1.TestGetMeContact())

	time.Sleep(time.Second)
}
