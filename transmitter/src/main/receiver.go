package main

import (
	"net/rpc"
	"github.com/RomanSaveljev/android-symbols/shared/src/shared"
	"fmt"
)

// Hides RPC details of the receiving end

type Receiver struct {
	client   *rpc.Client
	token string
}

// Registers object on the server side and creates a necessary file tree
// if this is a new file
func NewReceiver(fileName string, client *rpc.Client) (*Receiver, error) {
	var rx Receiver
	rx.client = client
	err := client.Call("Registrar.Create", &fileName, &rx.token)
	return &rx, err
}

// Retrieves next signature or returns error
func (this *Receiver) NextSignature() (shared.Signature, error) {
	var sig shared.Signature
	err := this.client.Call(fmt.Sprint(this.token, ".NextSignature"), nil, &sig)
	return sig, err
}

// Writes to the stream
func (this *Receiver) WriteStream(p []byte) error {
	err := this.client.Call(fmt.Sprint(this.token, ".WriteStream"), &p, nil)
	return err
}

func (this *Receiver) Close() error {
	return this.client.Call(fmt.Sprint(this.token, ".Close"), nil, nil)
}
