package main

import (

	"flag"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/timgranstrom/ptp-dfs"
	"time"
)

func main() {
	storeCommand := flag.NewFlagSet("store", flag.ExitOnError)
	catCommand := flag.NewFlagSet("cat", flag.ExitOnError)
	pinCommand := flag.NewFlagSet("pin", flag.ExitOnError)
	unpinCommand := flag.NewFlagSet("unpin", flag.ExitOnError)

	storeName := storeCommand.String("name", "", "name of the file")
	catName := catCommand.String("hash", "", "the hash of the file")
	pinName := pinCommand.String("hash", "", "the file-hash that should be pinned")
	unpinName := unpinCommand.String("hash", "", "the file-hash that should be pinned")

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

	default:
		log.Println("no commands")
	}

	if unpinCommand.Parsed() {
		fmt.Println("file that should be unpinned = " + *unpinName)
		//Call unpin function
	}

	if pinCommand.Parsed() {
		fmt.Println("file that should be pinned = " + *pinName)
		//Call pin function
	}

	if storeCommand.Parsed() {
		// call the store function and print the hash out
		count := 0
		for i := len(*storeName) - 1; i >= 0; i-- { // iteration through filePath to find only the filename
			if (string(*storeName)[i]) != 92 { // loop untill \ found
				count ++
			} else {
				break
			}
		}
		startIndex := len(*storeName) - count                 // startindexin orginal string
		fileName := string(*storeName)[startIndex:]          //slice out the FileName
		fmt.Println("final sliced string =     " + fileName) // print the final sliced string
		b, err := ioutil.ReadFile(*storeName)                 // Take out the content of the file in byte
		if err != nil {
			fmt.Print(err)
		}
		//fmt.Println("File stored = " + *storeName + "  " )
		fmt.Println(b)

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
}



