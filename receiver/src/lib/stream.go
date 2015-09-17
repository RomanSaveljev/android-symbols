package receiver

import (
	"os"
	"github.com/RomanSaveljev/android-symbols/shared/src/shared"
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

func (this *Stream) CloseAndOptimize(dummy int, signatures *shared.Signatures) error {
	err := this.stream.Close()
	// TODO: optimize and collect new signatures
	return err
}
