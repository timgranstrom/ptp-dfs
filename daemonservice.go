package ptp

import (
	"github.com/takama/daemon"
	"log"
	"os"
	"flag"
	"time"
)

const (
	defaultName = "dfs"
	defaultDescription = "This is a kademlia implementation of a distributed file system"
	defaultCommandPort = "7009"
	defaultKademliaPort = "8000"

)

// Service has embedded daemon
type DaemonService struct {
	daemon.Daemon
	dependencies  []string
}

func NewDaemonService() *DaemonService{
	service := &DaemonService{} //Create service
	newDaemon, err := daemon.New(defaultName,defaultDescription) //Create default daemon
	service.Daemon = newDaemon //Add default daemon to service

	if err != nil{
		log.Fatal("Error during daemon naming: ",err)
		os.Exit(1) // quit application
	}
	return service
}

func (ds *DaemonService) Install(name string, description string){
	ds.RenameDaemon(name)

	_, err := ds.Daemon.Install()
	if err != nil{
		log.Fatal("Error during daemon creation: ",err)
		os.Exit(1) // quit application
	}

	log.Println("Installed Daemon")
}

func (ds *DaemonService) Remove(name string){
	ds.RenameDaemon(name)

	_, err := ds.Daemon.Remove()
	if err != nil{
		log.Fatal("Error during install: ",err)
		os.Exit(1) // quit application
	}

	log.Println("Removed Daemon")
}

func (ds *DaemonService) RenameDaemon(name string){
	if name != "" {
		newDaemon, err := daemon.New(name, defaultDescription)
		if err != nil {
			log.Fatal("Error during daemon naming: ", err)
			os.Exit(1) // quit application
		}

		ds.Daemon = newDaemon
	}
}

func (ds *DaemonService) Start(name string, port string){
	kademliaListenPort := port

	if kademliaListenPort == "" {
		kademliaListenPort = defaultKademliaPort
	}

	kademliaInstance1 := NewKademlia(":8000")
	kademliaInstance1.routingTable.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"),":8000"))
	//kademliaInstance2 := NewKademlia(":8003")

	go kademliaInstance1.Run()
	//kademliaInstance2.Run()
	//time.Sleep(time.Second)
	go kademliaInstance1.TestSendMsg(kademliaInstance1.TestGetMeContact())

}


func (ds *DaemonService) HandleCommands(){

	// Commands
	installCommand := flag.NewFlagSet("remove", flag.ExitOnError)
	removeCommand := flag.NewFlagSet("install", flag.ExitOnError)
	startCommand := flag.NewFlagSet("start",flag.ExitOnError)

	//Install parameters
	nameInstallFlag := installCommand.String("name","","Name of daemon to install")
	infoInstallFlag := installCommand.String("info","","description of daemon to install")
	//Remove parameters
	nameRemoveFlag := removeCommand.String("name","","Name of daemon to remove")
	//Start parameters
	nameStartFlag := startCommand.String("name","","Name of daemon to start")
	kademliaListenPortStartFlag := startCommand.String("kademliaport","","Port number that the kademlia will listen on")

	switch os.Args[1] {
	case "install":
		installCommand.Parse(os.Args[2:]) //Parse the other commands into the install command
	case "remove":
		removeCommand.Parse(os.Args[2:]) //Parse the other commands into the remove command
	case "start":
		startCommand.Parse(os.Args[2:]) //Parse the other commands into the start command
	}

	if installCommand.Parsed() {
		ds.Install(*nameInstallFlag, *infoInstallFlag)
	}

	if removeCommand.Parsed() {
		ds.Remove(*nameRemoveFlag)
	}

	if startCommand.Parsed(){
		ds.Start(*nameStartFlag, *kademliaListenPortStartFlag)
	}
}