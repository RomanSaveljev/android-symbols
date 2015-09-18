package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
)

const APP_VERSION = "0.0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	//flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	//rest := flag.Args()

	var tr Transport
	log.Println("transport created")
	client := rpc.NewClient(&tr)
	log.Println("client created")
	client.Call("hello", true, false)
}
