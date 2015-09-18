package transmitter

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"sort"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
)

// Takes a list of signatures and produces a stream of literal bytes
// interspersed with matching block signatures. The client may stuff
// the data as they want, because compressor has internal data buffer
// able to hold one block of data.

type compressor struct {
	buffer     []byte
	enc        Encoder
	signatures *Signatures
}

func NewCompressor(signatures *Signatures, destination io.Writer) io.WriteCloser {
	var tx compressor
	tx.buffer = make([]byte, 0, receiver.CHUNK_SIZE)
	tx.enc = NewEncoder(destination)
	tx.signatures = signatures
	return &tx
}

func (this *compressor) writeFirst() error {
	_, err := this.enc.Write(this.buffer[:1])
	if err == nil {
		buffer := make([]byte, len(this.buffer) - 1, cap(this.buffer))
		copy(buffer, this.buffer[1:])
		this.buffer = buffer
	}
	return err
}

func (this *compressor) writeSignature(rolling string, signature string) error {
	log.Println("writeSignature")
	err := this.enc.WriteSignature(rolling, signature)
	if err == nil {
		this.buffer = make([]byte, 0, cap(this.buffer))
		log.Printf("new cap=%d", cap(this.buffer))
	}
	return err
}

func (this *compressor) writeOne(p byte) (err error) {
	log.Println("WriteOne")
	this.buffer = append(this.buffer, p)
	log.Printf("cap = %d len = %d", cap(this.buffer), len(this.buffer))
	if len(this.buffer) == cap(this.buffer) {
		log.Println("len reached cap and buf=%s", string(this.buffer))
		rolling := fmt.Sprintf("%08x", crc32.ChecksumIEEE(this.buffer))
		log.Printf("crc=%v", rolling)
		candidates := this.signatures.Get(rolling)
		log.Printf("candidates len=%d", len(candidates))
		if len(candidates) == 0 {
			log.Println("rolling checksum not found")
			err = this.writeFirst()
		} else {
			strong := fmt.Sprintf("%x", md5.Sum(this.buffer))
			log.Printf("strong=%s", strong)
			idx := sort.Search(len(candidates), func (i int) bool {return strong == candidates[i]})
			log.Printf("idx=%d", idx)
			if idx == len(candidates) {
				log.Println("strong signature not found")
				err = this.writeFirst()
			} else {
				err = this.writeSignature(rolling, candidates[idx])
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
	log.Printf("Close len=%d buffer=%s", len(this.buffer), this.buffer)
	_, err = this.enc.Write(this.buffer)
	if err == nil {
		err = this.enc.Close()
	}
	return err
}
