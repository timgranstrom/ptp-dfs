package main

import (
	"github.com/timgranstrom/ptp-dfs"
	"time"
)

func main() {
	kademliaInstance1 := ptp.NewKademlia(":8001")
	kademliaInstance2 := ptp.NewKademlia(":8002")
	kademliaInstance3 := ptp.NewKademlia(":8003")
	kademliaInstance3.TestAddContact(*kademliaInstance2.TestGetMeContact())
	kademliaInstance2.TestAddContact(*kademliaInstance1.TestGetMeContact())

	kademliaInstance1.Run()
	kademliaInstance2.Run()
	kademliaInstance3.Run()
	time.Sleep(time.Second)
	go kademliaInstance3.TestSendMsg(kademliaInstance1.TestGetMeContact())

	//kademliaInstance1.
	for{
		time.Sleep(time.Second)
	}
}


