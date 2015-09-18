package transmitter

import (
	"io"
	"log"
	"sort"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because compressor has internal data buffer
// able to hold one block of data.

type compressor struct {
	Chunk
	chunker  Chunker
	receiver Receiver
	buffer   []byte
}

func NewCompressor(chunker Chunker, rcv Receiver) io.WriteCloser {
	tx := compressor{chunker: chunker, receiver: rcv}
	tx.emptyBuffer()
	return &tx
}

func (this *compressor) emptyBuffer() {
	this.buffer = this.Data[:0]
}

func (this *compressor) isFull() bool {
	return len(this.buffer) == cap(this.buffer)
}

func (this *compressor) shiftData() {
	for i := 0; i < len(this.buffer)-1; i++ {
		this.Data[i] = this.Data[i+1]
	}
	this.buffer = this.Data[0 : len(this.buffer)-1]	
}

func (this *compressor) writeFirst() error {
	err := this.chunker.Write(this.buffer[0])
	if err == nil {
		this.shiftData()
	}
	return err
}

func (this *compressor) writeSignature(rolling string, signature string) error {
	log.Println("writeSignature")
	err := this.chunker.WriteSignature(rolling, signature)
	if err == nil {
		this.emptyBuffer()
		log.Printf("new cap=%d", cap(this.buffer))
	}
	return err
}

func (this *compressor) tryWriteSignature() (err error) {
	if signatures, err := this.receiver.Signatures(); err == nil {
		rolling := this.CountRolling()
		candidates := signatures.Get(rolling)
		if len(candidates) == 0 {
			err = this.writeFirst()
		} else {
			strong := this.CountStrong()
			idx := sort.Search(len(candidates), func(i int) bool { return strong == candidates[i] })
			if idx == len(candidates) {
				log.Println("strong signature not found")
				err = this.writeFirst()
			} else {
				err = this.writeSignature(rolling, candidates[idx])
			}
		}
	}
	return
}

func (this *compressor) writeOne(p byte) (err error) {
	log.Println("WriteOne")
	this.buffer = append(this.buffer, p)
	log.Printf("cap = %d len = %d", cap(this.buffer), len(this.buffer))
	if this.isFull() {
		if err = this.tryWriteSignature(); err != nil {
			this.buffer = this.buffer[0 : len(this.buffer)-1]
		}
	}
	return err
}

func (this *compressor) Write(p []byte) (n int, err error) {
	err = nil
	for n = 0; n < len(p) && err == nil; n++ {
		err = this.writeOne(p[n])
	}
	return n, err
}

func (this *compressor) Close() (err error) {
	log.Printf("Close len=%d buffer=%s", len(this.buffer), this.buffer)
	for len(this.buffer) != 0 {
		if err = this.chunker.Write(this.buffer[0]); err == nil {
			this.shiftData()
		}
	}
	return
}
