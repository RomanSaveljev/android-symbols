package main

import (
	"os"
)

type transport struct {
}

func (this *transport) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (this *transport) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (this *transport) Close() error {
	return os.Stdout.Close()
}
