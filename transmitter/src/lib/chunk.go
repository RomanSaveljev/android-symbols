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

func (this *Chunk) CountRolling() {
	this.Rolling = fmt.Sprintf("%08x", crc32.ChecksumIEEE(this.Data[:]))
}

func (this *Chunk) CountStrong() {
	this.Strong = fmt.Sprintf("%x", md5.Sum(this.Data[:]))
}