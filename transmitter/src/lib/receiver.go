package transmitter

import (
	"fmt"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"net/rpc"
	"io"
)

// Hides RPC details of the receiving end

type Receiver interface {
	io.WriteCloser
	NextSignature() (receiver.Signature, error)
	SaveChunk(chunk *receiver.Chunk) error
}

//go:generate $GOPATH/bin/mockgen -package transmitter -destination mock_receiver_test.go github.com/RomanSaveljev/android-symbols/transmitter/src/lib Receiver

type realReceiver struct {
	client *rpc.Client
	token  string
	stream string
}

// Registers object on the server side and creates a necessary file tree
// if this is a new file
func NewReceiver(fileName string, client *rpc.Client) (Receiver, error) {
	var rx realReceiver
	rx.client = client
	err := client.Call("Registrar.Create", &fileName, &rx.token)
	return &rx, err
}

// Retrieves next signature or returns error
func (this *realReceiver) NextSignature() (receiver.Signature, error) {
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
	return this.client.Call(fmt.Sprint(this.token, ".SaveChunk"), chunk, nil)
}
