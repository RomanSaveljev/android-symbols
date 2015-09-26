package receiver

import (
	_ "encoding/binary"
	"fmt"
	rxapp "github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/chunk"
	"github.com/RomanSaveljev/android-symbols/transmitter/signatures"
	"io"
	"net/rpc"
)

type Client interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
	Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call
}

type Receiver interface {
	io.WriteCloser
	SaveChunk(rolling uint32, strong, data []byte) error
	Signatures() (signatures.Signatures, error)
}

//go:generate $GOPATH/bin/mockgen -package mock -destination ../mock/mock_receiver.go github.com/RomanSaveljev/android-symbols/transmitter/receiver Receiver
//go:generate $GOPATH/bin/mockgen -package mock -destination ../mock/mock_client.go github.com/RomanSaveljev/android-symbols/transmitter/receiver Client

type realReceiver struct {
	client     Client
	token      string
	stream     string
	signatures signatures.Signatures
	call       *rpc.Call
	streamIndex int
}

// Registers object on the server side and creates a necessary file tree
// if this is a new file
func NewReceiver(fileName string, client Client) (Receiver, error) {
	var rx realReceiver
	rx.client = client
	err := client.Call("Synchronizer.StartFile", fileName, &rx.token)
	return &rx, err
}

func NewSecondaryReceiver(receiver Receiver, streamIndex int) Receiver {
	rcv := *(receiver.(*realReceiver))
	rcv.call = nil
	rcv.stream = ""
	rcv.streamIndex = streamIndex
	return &rcv
}

func (this *realReceiver) Signatures() (sigs signatures.Signatures, err error) {
	if this.signatures == nil {
		sigs = signatures.NewSignatures()
		for true {
			var sig rxapp.Signature
			if sig, err = this.nextSignature(); err == nil {
				var rolling uint32
				if rolling, err = chunk.RollingFromString(sig.Rolling); err != nil {
					continue
				}
				var strong []byte
				if strong, err = chunk.StrongFromString(sig.Strong); err != nil {
					continue
				}
				sigs.Add(rolling, strong)
			} else {
				break
			}
		}
		if err.Error() == io.EOF.Error() {
			this.signatures = sigs
			err = nil
		}
	} else {
		sigs = this.signatures
	}
	return
}

func (this *realReceiver) ensureRemoteIsFree() (err error) {
	if this.call != nil {
		done := <- this.call.Done
		err = done.Error
		this.call = nil
	} 
	return
}

// Retrieves next signature or returns error
func (this *realReceiver) nextSignature() (sig rxapp.Signature, err error) {
	if err = this.ensureRemoteIsFree(); err == nil {
		err = this.client.Call(fmt.Sprint(this.token, ".NextSignature"), 0, &sig)
	}
	return
}

// Writes to the stream
func (this *realReceiver) Write(p []byte) (n int, err error) {
	if err = this.ensureRemoteIsFree(); err == nil {
		if len(this.stream) == 0 {
			err = this.client.Call(fmt.Sprint(this.token, ".StartStream"), this.streamIndex, &this.stream)
		}
		if err == nil {
			this.call = this.client.Go(fmt.Sprint(this.stream, ".Write"), p, &n, nil)
		}
	}
	return
}

// Closes the stream
func (this *realReceiver) Close() (err error) {
	this.ensureRemoteIsFree()
	if len(this.stream) > 0 {
		err = this.client.Call(fmt.Sprint(this.stream, ".Close"), 0, nil)
		this.stream = ""
	}
	return
}

// Creates a new chunk
func (this *realReceiver) SaveChunk(rolling uint32, strong, data []byte) (err error) {
	var c rxapp.Chunk
	c.Rolling = chunk.RollingToString(rolling)
	c.Strong = chunk.StrongToString(strong)
	copy(c.Data[:], data)
	var sigs signatures.Signatures
	if sigs, err = this.Signatures(); err == nil {
		sigs.Add(rolling, strong)
	}
	if err = this.ensureRemoteIsFree(); err == nil {
		this.call = this.client.Go(fmt.Sprint(this.token, ".SaveChunk"), c, nil, nil)
	}
	return
}
