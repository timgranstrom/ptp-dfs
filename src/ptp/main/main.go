package main

import (
	"ptp"
	"time"
)

func main() {
	kademliaInstance1 := ptp.NewKademlia(":8000")
	go kademliaInstance1.Run()
	time.Sleep(time.Second)

	kademliaInstance2 := ptp.NewKademlia(":8005")
	go kademliaInstance2.TestSendMsg()
	for{
		time.Sleep(time.Second)
	}
}


