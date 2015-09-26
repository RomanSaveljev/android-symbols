package main

import (
	"flag"
	"fmt"
	//"github.com/RomanSaveljev/android-symbols/transmitter/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/gui"
	"github.com/RomanSaveljev/android-symbols/transmitter/processor"
	"github.com/RomanSaveljev/android-symbols/transmitter/receiver"
	"github.com/edsrzf/mmap-go"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"runtime/pprof"
	"strings"
)

const APP_VERSION = "0.0.1"

// RECEIVER=cmd.. PREFIX=... transmitter files...

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number")

func processOneFile(mm mmap.MMap, rcv receiver.Receiver, channel chan<- int) {
	processor.ProcessFileSync(mm, rcv, channel)
	close(channel)
}

func main() {
	profile := os.Getenv("CPU_PROFILE")
	if len(profile) > 0 {
		prof, _ := os.Create(profile)
		pprof.StartCPUProfile(prof)
	}

	log.Println("TX: starting")
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

	splitCmd := strings.Split(command, " ")
	tr, err := NewProcessTransport(exec.Command(splitCmd[0], splitCmd[1:]...))
	if err != nil {
		panic(fmt.Sprintf("Failed to create a transport: %v", err))
	}
	log.Println("TX: transport created")
	defer tr.Close()
	client := rpc.NewClient(tr)
	defer client.Close()
	var ui gui.Gui
	var channels []chan int
	for _, f := range rest {
		if file, err := os.Open(f); err == nil {
			defer file.Close()
			if mm, err := mmap.Map(file, mmap.RDONLY, 0); err == nil {
				defer mm.Unmap()
				ui.Total += uint64(len(mm))
				if rcv, err := receiver.NewReceiver(path.Join(prefix, f), client); err == nil {
					channels = append(channels, make(chan int))
					go processOneFile(mm[:len(mm)/2], rcv, channels[len(channels)-1])
					rcv := receiver.NewSecondaryReceiver(rcv, 1)
					channels = append(channels, make(chan int))
					go processOneFile(mm[len(mm)/2:], rcv, channels[len(channels)-1])
				}
			}
		}
	}

	inputs := make([]<-chan int, len(channels))
	for i, _ := range channels {
		inputs[i] = channels[i]
	}
	combined := gui.FanIn(inputs...)
	ui.Loop(combined)

	if len(profile) > 0 {
		pprof.StopCPUProfile()
	}
}
