package main

import (
	"os"
	"log"
)

type transport struct {
}

func (this *transport) Read(p []byte) (n int, err error) {
	log.Println("transport read")
	return os.Stdin.Read(p)
}

func (this *transport) Write(p []byte) (n int, err error) {
	log.Println("transport write: %v", p)
	return os.Stdout.Write(p)
}

func (this *transport) Close() error {
	return os.Stdout.Close()
}
