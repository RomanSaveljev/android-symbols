package transmitter

import (
	"fmt"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/signatures"
	"io"
)

type Client interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
}
//go:generate $GOPATH/bin/mockgen -package mock_transmitter -destination mock/mock_client.go github.com/RomanSaveljev/android-symbols/transmitter/src/lib Client

type Receiver interface {
	io.WriteCloser
	SaveChunk(chunk *receiver.Chunk) error
	Signatures() (*signatures.Signatures, error)
}

//go:generate $GOPATH/bin/mockgen -package mock_transmitter -destination mock/mock_receiver.go github.com/RomanSaveljev/android-symbols/transmitter/src/lib Receiver

type realReceiver struct {
	client Client
	token  string
	stream string
	signatures *signatures.Signatures
}

// Registers object on the server side and creates a necessary file tree
// if this is a new file
func NewReceiver(fileName string, client Client) (Receiver, error) {
	var rx realReceiver
	rx.client = client
	err := client.Call("Synchronizer.StartFile", fileName, &rx.token)
	return &rx, err
}

func (this *realReceiver) Signatures() (sigs *signatures.Signatures, err error) {
	if this.signatures == nil {
		sigs = signatures.NewSignatures()
		for true {
			if sig, err := this.nextSignature(); err == nil {
				sigs.Add(sig.Rolling, sig.Strong)
			} else {
				break
			}
		}
		if err == io.EOF {
			this.signatures = sigs
			err = nil
		}
	}
	return
}

// Retrieves next signature or returns error
func (this *realReceiver) nextSignature() (receiver.Signature, error) {
	var sig receiver.Signature
	err := this.client.Call(fmt.Sprint(this.token, ".NextSignature"), nil, &sig)
	return sig, err
}

// Writes to the stream
func (this *realReceiver) Write(p []byte) (n int, err error) {
	if len(this.stream) == 0 {
		err = this.client.Call(fmt.Sprint(this.token, ".StartStream"), 0, &this.stream)
	}
	if err == nil {
		err = this.client.Call(fmt.Sprint(this.stream, ".Write"), p, &n)
	}
	return
}

// Closes the stream
func (this *realReceiver) Close() (err error) {
	if len(this.stream) > 0 {
		err = this.client.Call(fmt.Sprint(this.stream, ".Close"), 0, nil)
		this.stream = ""
	}
	return
}

// Creates a new chunk
func (this *realReceiver) SaveChunk(chunk *receiver.Chunk) error {
	return this.client.Call(fmt.Sprint(this.token, ".SaveChunk"), *chunk, nil)
}
