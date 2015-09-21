package transmitter

import (
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	_ "github.com/RomanSaveljev/android-symbols/transmitter/src/lib/golang-ring"
	"io"
	_ "log"
	_ "fmt"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because Compressor has internal data buffer
// able to hold one block of data.

type Compressor struct {
	chunker  Chunker
	receiver Receiver
	//buffer   ring.Ring
	buffer []byte
	original []byte
}

func NewCompressor(chunker Chunker, rcv Receiver) io.WriteCloser {
	tx := Compressor{chunker: chunker, receiver: rcv}
	//tx.buffer.SetCapacity(receiver.CHUNK_SIZE)
	tx.buffer = make([]byte, 0, receiver.CHUNK_SIZE * 2)
	tx.original = tx.buffer
	//tx.emptyBuffer()
	return &tx
}

func (this *Compressor) emptyBuffer() {
	//this.buffer.Empty()
	// avoid reallocation at all cost!
	this.buffer = this.original
}

func (this *Compressor) isFull() bool {
	//return this.buffer.Length() == this.buffer.Capacity()
	return len(this.buffer) == receiver.CHUNK_SIZE
}

func (this *Compressor) writeFirst() error {
	//b := this.buffer.Dequeue()
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
		buffer, extra := this.buffer, []byte{} //this.buffer.Values()
		rolling := CountRolling(buffer, extra)
		candidates := signatures.Get(rolling)
		//log.Printf("rolling - %x", rolling)
		if candidates == nil {
			//log.Println("tryWriteSignature - no rolling candidates")
			err = this.writeFirst()
		} else {
			//log.Println("tryWriteSignature - rolling candidates:", candidates)
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
	//this.buffer.Enqueue(p)
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
	/*
	for this.buffer.Length() != 0 && err == nil {
		err = this.chunker.Write(this.buffer.Dequeue())
	}
	*/
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
