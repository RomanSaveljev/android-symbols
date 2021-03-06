package chunker

import (
	rxapp "github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/encoder"
	"github.com/RomanSaveljev/android-symbols/transmitter/receiver"
	"github.com/RomanSaveljev/android-symbols/transmitter/chunk"
	_ "log"
)

type Chunker interface {
	Flush() (err error)
	WriteSignature(rolling uint32, strong []byte) (err error)
	Write(b byte) (err error)
	Close() (err error)
}

//go:generate $GOPATH/bin/mockgen -package mock -destination ../mock/mock_chunker.go github.com/RomanSaveljev/android-symbols/transmitter/chunker Chunker

type realChunker struct {
	buffer   []byte
	encoder  encoder.Encoder
	receiver receiver.Receiver
}

func NewChunker(encoder encoder.Encoder, rcv receiver.Receiver) Chunker {
	var chunker = realChunker{encoder: encoder, receiver: rcv}
	chunker.buffer = make([]byte, 0, rxapp.CHUNK_SIZE)
	return &chunker
}

func (this *realChunker) emptyBuffer() {
	this.buffer = this.buffer[:0]
}

func (this *realChunker) isFull() bool {
	return len(this.buffer) == rxapp.CHUNK_SIZE
}

func (this *realChunker) Flush() (err error) {
	if this.isFull() {
		rolling := chunk.CountRolling(this.buffer)
		strong := chunk.CountStrong(this.buffer)
		if sigs, err := this.receiver.Signatures(); err == nil {
			group := sigs.Get(rolling)
			if group == nil || !group.Has(strong) {
				err = this.receiver.SaveChunk(rolling, strong, this.buffer)
			}
		}
		if err == nil {
			if err = this.writeSignature(rolling, strong); err == nil {
				this.emptyBuffer()
			}
		}
	} else {
		err = this.justFlush()
	}
	return
}

func (this *realChunker) justFlush() (err error) {
	if len(this.buffer) == 0 {
		return
	}
	var n int
	n, err = this.encoder.Write(this.buffer)
	if n > 0 && n < len(this.buffer) {
		copy(this.buffer[:len(this.buffer)-n], this.buffer[n:len(this.buffer)])
	}
	this.buffer = this.buffer[:len(this.buffer)-n]
	return
}

func (this *realChunker) WriteSignature(rolling uint32, strong []byte) (err error) {
	if err = this.justFlush(); err == nil {
		err = this.writeSignature(rolling, strong)
	}
	return
}

func (this *realChunker) writeSignature(rolling uint32, strong []byte) error {
	return this.encoder.WriteSignature(rolling, strong)
}

func (this *realChunker) Write(b byte) (err error) {
	this.buffer = append(this.buffer, b)
	if this.isFull() {
		if err = this.Flush(); err != nil {
			this.buffer = this.buffer[:len(this.buffer) - 1]
		}
	}
	return
}

func (this *realChunker) Close() (err error) {
	if err = this.Flush(); err == nil {
		err = this.encoder.Close()
	}
	return
}