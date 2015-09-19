package transmitter

import (
	"encoding/ascii85"
	"fmt"
	"io"
)

type Encoder interface {
	io.WriteCloser
	WriteSignature(rolling string, strong string) error
}

//go:generate $GOPATH/bin/mockgen -package mock_transmitter -destination mock/mock_encoder.go github.com/RomanSaveljev/android-symbols/transmitter/src/lib Encoder

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

type realEncoder struct {
	writer      *ascii85Writer
	destination io.Writer
}

func NewEncoder(destination io.Writer) Encoder {
	var enc realEncoder
	enc.writer = newAscii85Writer(destination)
	enc.destination = destination
	return &enc
}

func (this *realEncoder) Write(p []byte) (n int, err error) {
	return this.writer.Write(p)
}

func (this *realEncoder) WriteSignature(rolling string, strong string) error {
	err := this.writer.Close()
	this.writer = newAscii85Writer(this.destination)
	if err == nil {
		input := []byte(fmt.Sprintf("\t%s/%s\n", rolling, strong))
		_, err = this.destination.Write(input)
	}
	return err
}

func (this *realEncoder) Close() error {
	return this.writer.Close()
}
