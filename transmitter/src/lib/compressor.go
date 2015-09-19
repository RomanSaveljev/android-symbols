package transmitter

import (
	"io"
	"sort"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because Compressor has internal data buffer
// able to hold one block of data.

type Compressor struct {
	Chunk
	chunker  Chunker
	receiver Receiver
	buffer   []byte
}

func NewCompressor(chunker Chunker, rcv Receiver) io.WriteCloser {
	tx := Compressor{chunker: chunker, receiver: rcv}
	tx.emptyBuffer()
	return &tx
}

func (this *Compressor) emptyBuffer() {
	this.buffer = this.Data[:0]
}

func (this *Compressor) isFull() bool {
	return len(this.buffer) == cap(this.buffer)
}

func (this *Compressor) shiftData() {
	for i := 0; i < len(this.buffer)-1; i++ {
		this.Data[i] = this.Data[i+1]
	}
	this.buffer = this.Data[0 : len(this.buffer)-1]	
}

func (this *Compressor) writeFirst() error {
	err := this.chunker.Write(this.buffer[0])
	if err == nil {
		this.shiftData()
	}
	return err
}

func (this *Compressor) writeSignature(rolling string, signature string) error {
	err := this.chunker.WriteSignature(rolling, signature)
	if err == nil {
		this.emptyBuffer()
	}
	return err
}

func (this *Compressor) tryWriteSignature() (err error) {
	if signatures, err := this.receiver.Signatures(); err == nil {
		rolling := this.CountRolling()
		candidates := signatures.Get(rolling)
		if len(candidates) == 0 {
			err = this.writeFirst()
		} else {
			strong := this.CountStrong()
			idx := sort.Search(len(candidates), func(i int) bool { return strong == candidates[i] })
			if idx == len(candidates) {
				err = this.writeFirst()
			} else {
				err = this.writeSignature(rolling, candidates[idx])
			}
		}
	}
	return
}

func (this *Compressor) writeOne(p byte) (err error) {
	this.buffer = append(this.buffer, p)
	if this.isFull() {
		if err = this.tryWriteSignature(); err != nil {
			this.buffer = this.buffer[0 : len(this.buffer)-1]
		}
	}
	return err
}

func (this *Compressor) Write(p []byte) (n int, err error) {
	err = nil
	for n = 0; n < len(p) && err == nil; n++ {
		err = this.writeOne(p[n])
	}
	return n, err
}

func (this *Compressor) Close() (err error) {
	for len(this.buffer) != 0 && err == nil {
		if err = this.chunker.Write(this.buffer[0]); err == nil {
			this.shiftData()
		}
	}
	if err == nil {
		err = this.chunker.Close()
	}
	return
}
