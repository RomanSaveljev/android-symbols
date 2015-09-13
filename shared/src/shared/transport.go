package shared

import (
	"os"
	"log"
)

type Transport struct {
}

func (this *Transport) Read(p []byte) (n int, err error) {
	log.Println("transport read")
	return os.Stdin.Read(p)
}

func (this *Transport) Write(p []byte) (n int, err error) {
	log.Println("transport write: %v", p)
	return os.Stdout.Write(p)
}

func (this *Transport) Close() error {
	return os.Stdout.Close()
}
