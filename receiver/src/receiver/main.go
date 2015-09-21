package main

import (
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"log"
)

func main() {
	log.Println("RX: starting")
	var tr transport
	receiver.RunSynchronizerService(&tr)
}
