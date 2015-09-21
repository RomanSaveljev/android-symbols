package transmitter

import (
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"io"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because Compressor has internal data buffer
// able to hold one block of data.

type Compressor struct {
	chunker  Chunker
	receiver Receiver
	buffer []byte
	original []byte
	roller Roller
}

func NewCompressor(chunker Chunker, rcv Receiver) io.WriteCloser {
	tx := Compressor{chunker: chunker, receiver: rcv}
	tx.buffer = make([]byte, 0, receiver.CHUNK_SIZE * 2)
	tx.original = tx.buffer
	return &tx
}

func (this *Compressor) emptyBuffer() {
	this.buffer = this.original
	this.roller = Roller{}
}

func (this *Compressor) isFull() bool {
	return len(this.buffer) == receiver.CHUNK_SIZE
}

func (this *Compressor) writeFirst() error {
	b := this.buffer[0]
	this.buffer = this.buffer[1:]
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
		if this.roller.Calculated() {
			this.roller.Next(this.buffer[0], this.buffer[len(this.buffer) - 1])
		} else {
			this.roller.Calculate(this.buffer)
		}
		rolling := this.roller.Value()
		candidates := signatures.Get(rolling)
		if candidates == nil {
			err = this.writeFirst()
		} else {
			strong := CountStrong(this.buffer)
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
	this.buffer = append(this.buffer, p)
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
	for _, b := range this.buffer {
		if err = this.chunker.Write(b); err != nil {
			break
		}
	}
	if err == nil {
		err = this.chunker.Close()
	}
	return
}
