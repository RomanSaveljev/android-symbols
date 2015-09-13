package main

import (
	"os"
	"bufio"
	"io"
)

type transport struct {
	stdin io.Reader
	stdout io.Writer
}

func NewTransport() io.ReadWriteCloser {
	var tr transport
	tr.stdin = bufio.NewReader(os.Stdin)
	tr.stdout = bufio.NewWriter(os.Stdout)
	return &tr
}

func (this *transport) Read(p []byte) (n int, err error) {
	return this.stdin.Read(p)
}

func (this *transport) Write(p []byte) (n int, err error) {
	return this.stdout.Write(p)
}

func (this *transport) Close() error {
	return os.Stdout.Close()
}
