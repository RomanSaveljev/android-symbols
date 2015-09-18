package transmitter

import (
	"errors"
)

var ErrBufferIsFull = errors.New("Buffer is full")

type Chunker interface {
	Flush() (err error)
	WriteSignature(rolling string, strong string) (err error)
	Write(b byte) (err error)
	Close() (err error)
}

type realChunker struct {
	Chunk
	buffer   []byte
	encoder  Encoder
	receiver Receiver
}

func NewChunker(encoder Encoder, receiver Receiver) Chunker {
	var chunker = realChunker{encoder: encoder, receiver: receiver}
	chunker.buffer = chunker.Data[:0]
	return &chunker
}

func (this *realChunker) emptyBuffer() {
	this.buffer = this.Data[:0]
}

func (this *realChunker) isFull() bool {
	return len(this.buffer) == cap(this.buffer)
}

func (this *realChunker) Flush() (err error) {
	if this.isFull() {
		this.CountRolling()
		this.CountStrong()
		if err = this.receiver.SaveChunk(&this.Chunk.Chunk); err == nil {
			if err = this.writeSignature(this.Chunk.Rolling, this.Chunk.Strong); err == nil {
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

func (this *realChunker) WriteSignature(rolling string, strong string) (err error) {
	if err = this.justFlush(); err == nil {
		err = this.writeSignature(rolling, strong)
	}
	return
}

func (this *realChunker) writeSignature(rolling string, strong string) error {
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