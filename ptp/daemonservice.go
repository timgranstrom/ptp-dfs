package ptp

import (
	"github.com/takama/daemon"
	"log"
	"os"
	"fmt"
	"net"
	"path/filepath"
	"bufio"
	"strings"
	"io/ioutil"
	"encoding/hex"
)

const (
	DefaultDaemonName   = "dfs"
	DefaultDaemonDescription  = "This is a kademlia implementation of a distributed file system"
	DefaultDaemonCommandPort  = ":7009"
	DefaultKademliaPort  = ":8050"
)

// Service has embedded daemon
type DaemonService struct {
	daemon.Daemon
	dependencies  []string
	kademliaNode *Kademlia
}

func NewDaemonService() *DaemonService{
	service := &DaemonService{}                                        //Create service
	newDaemon, err := daemon.New(DefaultDaemonName,DefaultDaemonDescription) //Create default daemon
	service.Daemon = newDaemon                                         //Add default daemon to service
	service.kademliaNode = NewKademlia(DefaultKademliaPort,nil) //Created kademlia node
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

func (ds *DaemonService) RunKademlia(){
	ds.kademliaNode.Run()
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
	size, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	receivedMsg := string(buf[0:size]) //get message based on it's size
	log.Println("Received msg: ",receivedMsg)
	resultMessage := ds.ParseReceivedMessage(receivedMsg) //parse the message and handle it
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

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	conn.Write([]byte(message)) //Send message to daemon service
	// listen for reply
	size, error := bufio.NewReader(conn).Read(buf)
	if error != nil{
		log.Println("COULDN'T RECIEVE RESPONSE MESSAGE FROM DAEMON, Error:",error)
	}

	log.Print("Response from daemon: "+string(buf[0:size]))
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
	if len(msgArgs) > 0{ //Message contains something
		switch msgArgs[0] {
		case "store":
			return ds.Store(msgArgs[1])
		case "cat":
			return ds.Cat(msgArgs[1])
		case "pin":
			return ds.Pin(msgArgs[1],true)
		case "unpin":
			return ds.Pin(msgArgs[1],false)
		case "me":
			me := ds.kademliaNode.GetMe() //Get my own Kademlia contact
			return me.ID.String() //Return my kademlia ID
		}
	}
	return "No Valid Response"
}

func (ds *DaemonService) Cat(key string) string{
	decodedKey,err := hex.DecodeString(key)
	if err != nil{
		//log.Fatal("Key decoding error: ",err)
		return "Key decoding error: "+err.Error()
	}

	data,isfound := ds.kademliaNode.network.store.RetrieveData(decodedKey)

	if isfound{
		return string(data)
	}else{
		dataFoundWithLookup := ds.kademliaNode.LookupData(key)
		if dataFoundWithLookup{
			data,_ := ds.kademliaNode.network.store.RetrieveData(decodedKey)
			return string(data)
		}
	}
	return "Could not find data with key"
}

func (ds *DaemonService)Store(filepath string) string{
	// call the store function and print the hash out
	fileName,absPath := ds.ParseFilePathCommand(filepath)
	//fmt.Println("File name: ",fileName,", File Path:",absPath) // print the file name and it's path
	b, err := ioutil.ReadFile(absPath)                 // Take out the content of the file in byte
	if err != nil {
		//fmt.Print(err)
		return "Could not read file: "+err.Error()
	}
	key := ds.kademliaNode.Store(fileName,b)
	return key
}

func (ds *DaemonService) Pin(key string,isPinned bool) string{
	decodedKey,err := hex.DecodeString(key)
	if err != nil{
		//log.Fatal("Key decoding error: ",err)
		return "Key decoding error: "+err.Error()
	}

	_,isfound := ds.kademliaNode.network.store.RetrieveData(decodedKey) //Check if any data actually is stored for key

	if !isfound{
		return "Could not find data to pin with key"
	}

	err = ds.kademliaNode.network.store.SetPin(decodedKey,isPinned)
	if err != nil{
		//log.Fatal("Could not PIN for key, error: ",err)
		return "Key decoding error: "+err.Error()
	}
	if isPinned {
		return "Pin Successful"
	} else{
		return "Unpin Successful"
	}
}
