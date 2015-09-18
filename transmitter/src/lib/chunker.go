package transmitter

import (
	"errors"
)

var ErrBufferIsFull = errors.New("Buffer is full")

type Chunker struct {
	Chunk
	buffer   []byte
	encoder  Encoder
	receiver Receiver
}

func NewChunker(encoder Encoder, receiver Receiver) *Chunker {
	var chunker = Chunker{encoder: encoder, receiver: receiver}
	chunker.buffer = chunker.Data[:0]
	return &chunker
}

func (this *Chunker) emptyBuffer() {
	this.buffer = this.Data[:0]
}

func (this *Chunker) isFull() bool {
	return len(this.buffer) == cap(this.buffer)
}

func (this *Chunker) Flush() (err error) {
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

func (this *Chunker) justFlush() (err error) {
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

func (this *Chunker) WriteSignature(rolling string, strong string) (err error) {
	if err = this.justFlush(); err == nil {
		err = this.writeSignature(rolling, strong)
	}
	return
}

func (this *Chunker) writeSignature(rolling string, strong string) error {
	return this.encoder.WriteSignature(rolling, strong)
}

func (this *Chunker) Write(b byte) (err error) {
	if this.isFull() {
		// have to flush it manually
		err = ErrBufferIsFull
	} else {
		this.buffer = append(this.buffer, b)
		if this.isFull() {
			err = this.Flush()
		}
	}
	return
}

func (this *Chunker) Close() (err error) {
	if err = this.Flush(); err == nil {
		err = this.encoder.Close()
	}
	return
}