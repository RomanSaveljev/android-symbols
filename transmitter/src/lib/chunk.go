package transmitter

import (
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"hash/crc32"
	"crypto/md5"
	"fmt"
)

type Chunk struct {
	receiver.Chunk
}

func (this *Chunk) CountRolling() string {
	this.Rolling = fmt.Sprintf("%08x", crc32.ChecksumIEEE(this.Data[:]))
	return this.Rolling
}

func (this *Chunk) CountStrong() string {
	this.Strong = fmt.Sprintf("%x", md5.Sum(this.Data[:]))
	return this.Strong
}