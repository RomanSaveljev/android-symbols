package receiver

import (
	"os"
	"io"
)

type Stream struct {
	stream io.WriteCloser
}

func NewStream(filePath string) (*Stream, error) {
	var err error
	var stream Stream
	stream.stream, err = os.Create(filePath)
	return &stream, err
}

func (this *Stream) Write(data []byte, n *int) (err error) {
	*n, err = this.stream.Write(data)
	return err
}

func (this *Stream) Close(dummy int, nothing *int) error {
	err := this.stream.Close()
	return err
}
