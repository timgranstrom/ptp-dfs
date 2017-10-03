package daemons

import (
	"github.com/takama/daemon"
	"log"
	"os"
	"fmt"
	"net"
	"path/filepath"
	"bufio"
	"strings"
)

const (
	DefaultDaemonName   = "dfs"
	DefaultDaemonDescription  = "This is a kademlia implementation of a distributed file system"
	DefaultDaemonCommandPort  = ":7009"
	DefaultKademliaPort  = ":8000"

)

// Service has embedded daemon
type DaemonService struct {
	daemon.Daemon
	dependencies  []string
	//kademliaNode Kademlia
}

func NewDaemonService() *DaemonService{
	service := &DaemonService{}                                        //Create service
	newDaemon, err := daemon.New(DefaultDaemonName,DefaultDaemonDescription) //Create default daemon
	service.Daemon = newDaemon                                         //Add default daemon to service
	//service.kademliaNode = ptp.NewKademlia(DefaultKademliaPort,nil)
	if err != nil{
		log.Fatal("Error during daemon naming: ",err)
		os.Exit(1) // quit application
	}
	return service
}

func (ds *DaemonService) Install(){
	_, err := ds.Daemon.Install()
	if err != nil{
		log.Fatal("Error during daemon creation: ",err)
		os.Exit(1) // quit application
	}

	log.Println("Installed Daemon")
}

func (ds *DaemonService) Remove(){
	_, err := ds.Daemon.Remove()
	if err != nil{
		log.Fatal("Error during install: ",err)
		os.Exit(1) // quit application
	}

	log.Println("Removed Daemon")
}

func (ds *DaemonService) Start(){
	_, err := ds.Daemon.Start()
	if err != nil{
		log.Fatal("Error during daemon startup: ",err)
		os.Exit(1) // quit application
	}

	log.Println("Started Daemon")
}

func (ds *DaemonService) Stop(){
	_, err := ds.Daemon.Stop()
	if err != nil{
		log.Fatal("Error during daemon stop: ",err)
		os.Exit(1) // quit application
	}

	log.Println("Stopped Daemon")
}

func (ds *DaemonService) Status(){
	stats, err := ds.Daemon.Status()
	if err != nil{
		log.Fatal("Error while fetching daemon status: ",err)
		os.Exit(1) // quit application
	}
	log.Println(stats)
}

func (ds *DaemonService) RunDaemonCommandListener(){
	// Listen for incoming connections.
	l, err := net.Listen("tcp", DefaultDaemonCommandPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the command listener when the application closes.
	defer l.Close()

	fmt.Println("Listening for daemon commands on " + l.Addr().String())

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go ds.handleRequest(conn)
	}
}

// Handles incoming requests.
func (ds *DaemonService) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	log.Println("Received msg: ",string(buf))
	resultMessage := ds.ParseReceivedMessage(string(buf)) //parse the message and handle it

	// Send a response back to sender
	conn.Write([]byte(resultMessage))
	// Close the connection when finished
	conn.Close()
}

// Handles sending requests.
func (ds *DaemonService) SendRequest(message string) {
	conn, error := net.Dial("tcp", DefaultDaemonCommandPort)
	defer conn.Close()
	if error != nil{
		log.Println("COULDN'T SEND MESSAGE TO DAEMON, Error:",error)
	}
	conn.Write([]byte(message)) //Send message to daemon service
	// listen for reply
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Response from daemon: "+response)
}

func (ds *DaemonService) ParseFilePathCommand(relativeFilePath string) (fileName string, absoluteFilePath string){
	absFilePath,error := filepath.Abs(relativeFilePath) //Get absolute filepath representation
	if error != nil{
		log.Fatal("File path parsing error: ",error)
	}
	count := 0
	for i := len(absFilePath) - 1; i >= 0; i-- { // iteration through filePath to find only the filename
		if (string(absFilePath)[i]) != 92 { // loop until '\' found
			count ++
		} else {
			break
		}
	}
	startIndex := len(absFilePath) - count // start indexing original string
	fileName = string(absFilePath)[startIndex:] //slice out the FileName
	fileName = string(fileName)
	return fileName,absFilePath
}

/**
Parse the message and handle it accordingly
 */
func (ds *DaemonService) ParseReceivedMessage(message string) string{
	msgArgs := strings.Fields(message) //Split command message into fields
	log.Println("ARGUMENT LENGTH: ",len(msgArgs))
	if len(msgArgs) > 2{ //Message contains something
		switch msgArgs[1] {
		case "store":
			//storeCommand.Parse(os.Args[2:])
		case "cat":
			//catCommand.Parse(os.Args[2:])

			//log.Println(*catName)

		case "pin":
			//pinCommand.Parse(os.Args[2:])

		case "unpin":
			//unpinCommand.Parse(os.Args[2:])

		case "start":
			//daemonCommand.Parse(os.Args[2:])
		}
	}
	return "nothing done yet"
}
