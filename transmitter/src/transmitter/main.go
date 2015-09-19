package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"os/exec"
	"os"
	"strings"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib"
	"path"
)

const APP_VERSION = "0.0.1"

// RECEIVER=cmd.. transmitter files...
 
// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number")

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		os.Exit(0)
	}

	rest := flag.Args()

	command := os.Getenv("RECEIVER")
	if len(command) == 0 {
		panic("RECEIVER environment variable must tell receiver command")
	}
	
	prefix := os.Getenv("PREFIX")
	
	splitCmd := strings.Split(" ", command)
	tr, err := NewProcessTransport(exec.Command(splitCmd[0], splitCmd[1:]...))
	if err != nil {
		panic("Failed to create a transport")
	}
	defer tr.Close()
	client := rpc.NewClient(tr)
	defer client.Close()
	for _, f := range rest {
		if file, err := os.Open(f); err == nil {
			rcv, _ := transmitter.NewReceiver(path.Join(prefix, f), client)
			transmitter.ProcessFileSync(file, rcv)
			file.Close()
		}		
	}
}
