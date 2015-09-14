package receiver

import (
	"os"
	"github.com/RomanSaveljev/android-symbols/shared/src/shared"
)

type Stream struct {
	file *os.File
}

func NewStream(filePath string) (*Stream, error) {
	var err error
	var stream Stream
	stream.file, err = os.Create(filePath)
	return &stream, err
}

func (this *Stream) Write(data []byte, n *int) (err error) {
	*n, err = this.file.Write(data)
	return err
}

func (this *Stream) CloseAndOptimize(dummy int, signatures *shared.Signatures) error {
	err := this.file.Close()
	// TODO: optimize and collect new signatures
	return err
}
