package transmitter

import (
	"encoding/ascii85"
	"fmt"
	"io"
)

type ascii85Writer struct {
	writer       io.WriteCloser
	destination  io.Writer
	needLineFeed bool
}

func newAscii85Writer(destination io.Writer) *ascii85Writer {
	var writer ascii85Writer
	writer.writer = ascii85.NewEncoder(destination)
	writer.destination = destination
	writer.needLineFeed = false
	return &writer
}

func (this *ascii85Writer) Write(p []byte) (int, error) {
	n, err := this.writer.Write(p)
	if err == nil {
		this.needLineFeed = true
	}
	return n, err
}

func (this *ascii85Writer) Close() error {
	err := this.writer.Close()
	if err == nil && this.needLineFeed {
		_, err = this.destination.Write([]byte("\n"))
	}
	return err
}

type encoder struct {
	writer      *ascii85Writer
	destination io.Writer
}

func newEncoder(destination io.Writer) *encoder {
	var enc encoder
	enc.writer = newAscii85Writer(destination)
	enc.destination = destination
	return &enc
}

func (this *encoder) Write(p byte) error {
	var err error
	_, err = this.writer.Write([]byte{p})
	return err
}

func (this *encoder) WriteSignature(rolling uint32, strong string) error {
	err := this.writer.Close()
	this.writer = newAscii85Writer(this.destination)
	if err == nil {
		input := []byte(fmt.Sprintf("\t%08x/%s\n", rolling, strong))
		_, err = this.destination.Write(input)
	}
	return err
}

func (this *encoder) Close() error {
	return this.writer.Close()
}
