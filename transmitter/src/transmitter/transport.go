package main

import (
	"os/exec"
	"io"
)

type Transport struct {
	write io.WriteCloser
	read io.ReadCloser
	cmd *exec.Cmd
}

func NewProcessTransport(cmd *exec.Cmd) (*Transport, error) {
	var tr Transport
	var err error = nil
	tr.write, err = cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	tr.read, err = cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		tr.read.Close()
		tr.write.Close()
		return nil, err
	}
	tr.cmd = cmd
	return &tr, err
}

func (this *Transport) Read(p []byte) (n int, err error) {
	return this.read.Read(p)
}

func (this *Transport) Write(p []byte) (n int, err error) {
	return this.write.Write(p)
}

func (this *Transport) Close() error {
	this.write.Close()
	this.read.Close()
	return this.cmd.Wait()
}
