package main

import (
	"flag"
	"os"
	"log"
	"fmt"
	"github.com/timgranstrom/ptp-dfs/ptp"
)

func main() {
	daemonService := ptp.NewDaemonService()

	//Head commands
	storeCommand := flag.NewFlagSet("store", flag.ExitOnError)
	catCommand := flag.NewFlagSet("cat", flag.ExitOnError)
	pinCommand := flag.NewFlagSet("pin", flag.ExitOnError)
	unpinCommand := flag.NewFlagSet("unpin", flag.ExitOnError)
	meCommand := flag.NewFlagSet("me", flag.ExitOnError)
	daemonCommand := flag.NewFlagSet("daemon", flag.ExitOnError)

	//Sub commands
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

		case "pin":
			pinCommand.Parse(os.Args[2:])

		case "unpin":
			unpinCommand.Parse(os.Args[2:])

		case "daemon":
			daemonCommand.Parse(os.Args[2:])

		case "me":
			meCommand.Parse(os.Args[2:])
		}
	} else{
		log.Println("no commands")
		log.Println("[STARTING DAEMON COMMAND LISTENER AND KADEMLIA NODE]")
		go daemonService.RunKademlia() //Start up kademlia in parallell
		daemonService.RunDaemonCommandListener() //run the daemon command listener
	}

	if unpinCommand.Parsed() {
		fmt.Println("file that should be unpinned = " + *unpinName)
		unpinMessage := "unpin "+*unpinName
		daemonService.SendRequest(unpinMessage)
	}

	if pinCommand.Parsed() {
		pinMessage := "pin "+*pinName
		daemonService.SendRequest(pinMessage)

	}

	if storeCommand.Parsed() {
		storeMessage := "store "+*storeName
		daemonService.SendRequest(storeMessage)
	}

	if catCommand.Parsed() {
		catMessage := "cat "+*catName
		daemonService.SendRequest(catMessage)
	}

	if meCommand.Parsed() {
		meMessage := "me"
		daemonService.SendRequest(meMessage)
	}

	if daemonCommand.Parsed() {
		if *daemonInstall{
			daemonService.Install()
		}else if *daemonRemove{
			daemonService.Remove()

		}else if *daemonStart{
			daemonService.Start()

		}else if *daemonStop{
			daemonService.Stop()

		}else if *daemonStatus{
			daemonService.Status()
		}
	}
}



