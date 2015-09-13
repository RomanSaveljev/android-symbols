package transmitter

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io"
	"sort"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because compressor has internal data buffer
// able to hold one block of data.

type compressor struct {
	buffer     []byte
	enc        *encoder
	signatures Signatures
}

func NewTransmitter(signatures Signatures, bufferSize uint, destination io.Writer) io.WriteCloser {
	var tx compressor
	tx.buffer = make([]byte, bufferSize)
	tx.enc = newEncoder(destination)
	tx.signatures = signatures
	return &tx
}

func (this *compressor) writeFirstAndAppend(p byte) error {
	err := this.enc.Write(this.buffer[0])
	if err == nil {
		this.buffer = append(this.buffer[1:], p)
	}
	return err
}

func (this *compressor) writeSignatureAndAppend(rolling uint32, signature string, p byte) error {
	err := this.enc.WriteSignature(rolling, signature)
	if err == nil {
		this.buffer = []byte{p}
	}
	return err
}

func (this *compressor) writeOne(p byte) (err error) {
	this.buffer = append(this.buffer, p)
	if len(this.buffer) == cap(this.buffer) {
		rolling := crc32.ChecksumIEEE(this.buffer)
		candidates := this.signatures.Get(rolling)
		if len(candidates) == 0 {
			err = this.writeFirstAndAppend(p)
		} else {
			strong := fmt.Sprintf("%x", md5.Sum(this.buffer))
			idx := sort.SearchStrings(candidates, strong)
			if idx == -1 {
				err = this.writeFirstAndAppend(p)
			} else {
				err = this.writeSignatureAndAppend(rolling, candidates[idx], p)
			}
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

func (this *compressor) Close() error {
	var err error
	for i := 0; i < len(this.buffer) && err == nil; i++ {
		err = this.enc.Write(this.buffer[i])
	}
	if err == nil {
		err = this.enc.Close()
	}
	return err
}
