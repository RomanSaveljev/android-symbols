package transmitter

import (
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/golang-ring"
	"io"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because Compressor has internal data buffer
// able to hold one block of data.

type Compressor struct {
	chunker  Chunker
	receiver Receiver
	buffer   ring.Ring
}

func NewCompressor(chunker Chunker, rcv Receiver) io.WriteCloser {
	tx := Compressor{chunker: chunker, receiver: rcv}
	tx.buffer.SetCapacity(receiver.CHUNK_SIZE)
	tx.emptyBuffer()
	return &tx
}

func (this *Compressor) emptyBuffer() {
	this.buffer.Empty()
}

func (this *Compressor) isFull() bool {
	return this.buffer.Length() == this.buffer.Capacity()
}

func (this *Compressor) writeFirst() error {
	b := this.buffer.Dequeue()
	err := this.chunker.Write(b)
	if err != nil {
		panic("TODO: implement unshifting data to ring buffer")
	}
	return err
}

func (this *Compressor) writeSignature(rolling []byte, strong []byte) error {
	err := this.chunker.WriteSignature(rolling, strong)
	if err == nil {
		this.emptyBuffer()
	}
	return err
}

func (this *Compressor) tryWriteSignature() (err error) {
	if signatures, err := this.receiver.Signatures(); err == nil {
		buffer, extra := this.buffer.Values()
		rolling := CountRolling(buffer, extra)
		candidates := signatures.Get(rolling)
		if candidates == nil {
			err = this.writeFirst()
		} else {
			strong := CountStrong(buffer, extra)
			if candidates.Has(strong) {
				err = this.writeSignature(rolling, strong)
			} else {
				err = this.writeFirst()
			}
		}
	}
	return
}

func (this *Compressor) writeOne(p byte) (err error) {
	this.buffer.Enqueue(p)
	if this.isFull() {
		if err = this.tryWriteSignature(); err != nil {
			panic("TODO: implement undo enqueue for ring buffer")
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
	for this.buffer.Length() != 0 && err == nil {
		err = this.chunker.Write(this.buffer.Dequeue())
	}
	if err == nil {
		err = this.chunker.Close()
	}
	return
}
