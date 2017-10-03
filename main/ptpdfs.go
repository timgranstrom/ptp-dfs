package main

import (
	"flag"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/timgranstrom/ptp-dfs"
	"time"
	"github.com/timgranstrom/ptp-dfs/daemons"
)

func main() {
	daemonService := daemons.NewDaemonService()

	storeCommand := flag.NewFlagSet("store", flag.ExitOnError)
	catCommand := flag.NewFlagSet("cat", flag.ExitOnError)
	pinCommand := flag.NewFlagSet("pin", flag.ExitOnError)
	unpinCommand := flag.NewFlagSet("unpin", flag.ExitOnError)
	daemonCommand := flag.NewFlagSet("daemon", flag.ExitOnError)


	storeName := storeCommand.String("name", "", "name of the file")
	catName := catCommand.String("hash", "", "the hash of the file")
	pinName := pinCommand.String("hash", "", "the file-hash that should be pinned")
	unpinName := unpinCommand.String("hash", "", "the file-hash that should be pinned")
	daemonInstall := daemonCommand.Bool("install", false, "install the daemon")
	daemonRemove := daemonCommand.Bool("remove", false, "remove the daemon")
	daemonStart := daemonCommand.Bool("start", false, "start the daemon")
	daemonStop := daemonCommand.Bool("stop", false, "stop the daemon")
	daemonStatus := daemonCommand.Bool("status", false, "status of the daemon")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "store":
			storeCommand.Parse(os.Args[2:])

		case "cat":
			catCommand.Parse(os.Args[2:])

			log.Println(*catName)

		case "pin":
			pinCommand.Parse(os.Args[2:])

		case "unpin":
			unpinCommand.Parse(os.Args[2:])

		case "daemon":
			daemonCommand.Parse(os.Args[2:])
		}
	} else{
		log.Println("no commands")
		log.Println("[STARTING DAEMON COMMAND LISTENER]")
		daemonService.RunDaemonCommandListener()
	}

	if unpinCommand.Parsed() {
		fmt.Println("file that should be unpinned = " + *unpinName)
		unpinMessage := "unpin "+*unpinName
		daemonService.SendRequest(unpinMessage)
		//Call unpin function
	}

	if pinCommand.Parsed() {
		fmt.Println("file that should be pinned = " + *pinName)
		//Call pin function
	}

	if storeCommand.Parsed() {
		// call the store function and print the hash out
		fileName,absPath := daemonService.ParseFilePathCommand(*storeName)
		fmt.Println("File name: ",fileName,", File Path:",absPath) // print the file name and it's path
		b, err := ioutil.ReadFile(*storeName)                 // Take out the content of the file in byte
		if err != nil {
			fmt.Print(err)
		}
		//fmt.Println("File stored = " + *storeName + "  " )
		fmt.Println(string(b))

		nodeTest := ptp.NewKademlia(":8001",nil)
		go nodeTest.Run()
		time.Sleep(time.Second)
		nodeTest.Store(fileName,b)
		time.Sleep(time.Second)
		// log.Println(kademlia.store(subString,b))
	}

	if catCommand.Parsed() {
		// find the path to the data that should be printed out
		//fileName =: Kademlia.findData(*catName)
		b, err := ioutil.ReadFile(*catName) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println(b)
		str := string(b) // convert content to a string
		fmt.Println(str)
	}

	if daemonCommand.Parsed() {
		if *daemonInstall{
			//log.Println("SUCCESS INSTALL TEST")
			daemonService.Install()
		}else if *daemonRemove{
			//log.Println("SUCCESS REMOVE TEST")
			daemonService.Remove()

		}else if *daemonStart{
			//log.Println("SUCCESS START TEST")
			daemonService.Start()

		}else if *daemonStop{
			//log.Println("SUCCESS STOP TEST")
			daemonService.Stop()

		}else if *daemonStatus{
			//log.Println("SUCCESS STATUS TEST")
			daemonService.Status()
		}
	}
}



