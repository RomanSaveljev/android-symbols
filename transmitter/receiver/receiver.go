package receiver

import (
	"encoding/hex"
	"fmt"
	rxapp "github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/signatures"
	"io"
)

type Client interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
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
}

// Registers object on the server side and creates a necessary file tree
// if this is a new file
func NewReceiver(fileName string, client Client) (Receiver, error) {
	var rx realReceiver
	rx.client = client
	err := client.Call("Synchronizer.StartFile", fileName, &rx.token)
	return &rx, err
}

func hexToDec(b byte) uint32 {
	/*
	if b < '0' {
		panic("Not a hex digit")
	} else if b > 'f' {
		panic("Not a hex digit")
	} else if b <= '9' {
		return uint32(b - '0')
	} else if b >= 'a' {
		return uint32(b - 'a' + 0xa)
	} else {
		panic("Not a hex digit")
	}
	panic("Should not go here")
	*/
	// The trivial "b - '0'" implementation is on par
	// with the below unfolded bubble search according
	// to the benchmark
	if b < 'a' {
		if b < '6' {
			if b < '3' {
				switch b {
				case '0':
					return 0
				case '1':
					return 1
				case '2':
					return 2
				}
			} else {
				switch b {
				case '3':
					return 3
				case '4':
					return 4
				case '5':
					return 5
				}
			}
		} else {
			switch b {
			case '6':
				return 6
			case '7':
				return 7
			case '8':
				return 8
			case '9':
				return 9
			}
		}
	} else {
		if b < 'd' {
			switch b {
			case 'a':
				return 10
			case 'b':
				return 11
			case 'c':
				return 12
			}
		} else {
			switch b {
			case 'd':
				return 13
			case 'e':
				return 14
			case 'f':
				return 15
			}
		}
	}
	panic("Not a hex digit passed")
}

func extractRolling(s string) (ret uint32) {
	if len(s) != 8 {
		panic("Unsupported string length")
	}
	ret = hexToDec(s[0])<<28 |
		hexToDec(s[1])<<24 |
		hexToDec(s[2])<<20 |
		hexToDec(s[3])<<16 |
		hexToDec(s[4])<<12 |
		hexToDec(s[5])<<8 |
		hexToDec(s[6])<<4 |
		hexToDec(s[7])
	return
}

func (this *realReceiver) Signatures() (sigs signatures.Signatures, err error) {
	if this.signatures == nil {
		sigs := signatures.NewSignatures()
		for true {
			var sig rxapp.Signature
			if sig, err = this.nextSignature(); err == nil {
				rolling := extractRolling(sig.Rolling)
				var strong []byte
				if strong, err = hex.DecodeString(sig.Strong); err != nil {
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

// Retrieves next signature or returns error
func (this *realReceiver) nextSignature() (rxapp.Signature, error) {
	var sig rxapp.Signature
	err := this.client.Call(fmt.Sprint(this.token, ".NextSignature"), 0, &sig)
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
func (this *realReceiver) SaveChunk(rolling uint32, strong, data []byte) error {
	var chunk rxapp.Chunk
	chunk.Rolling = fmt.Sprintf("%x", rolling)
	chunk.Strong = hex.EncodeToString(strong)
	copy(chunk.Data[:], data)
	return this.client.Call(fmt.Sprint(this.token, ".SaveChunk"), chunk, nil)
}
