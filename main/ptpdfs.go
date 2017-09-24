package main

import (
	"github.com/timgranstrom/ptp-dfs"
	"os"
	"time"
	"fmt"
)

/*func main() {
	/*kademliaInstance1 := ptp.NewKademlia(":8001")
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
	}*/

	/*var ip = flag.String("ip","","Ip address")
	var port = flag.String("p","","Port")
	flag.Parse()

	if *port == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println(*ip)*/


	/*srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)

	for {
		time.Sleep(time.Second)
		fmt.Println(status)
	}
}*/

func main() {

	ds := ptp.NewDaemonService()
	if len(os.Args) > 1 {
		ds.HandleCommands()
	}/*else{
		ds.Start("","") //run on default
		log.Println("RAN!")
	}*/
	/*for{
		time.Sleep(time.Second)
		fmt.Printf("RUNNING")
	}*/
}




