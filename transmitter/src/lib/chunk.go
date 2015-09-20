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

func GetRolling(buffer []byte) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE(buffer))
}

func (this *Chunk) CountRolling() string {
	this.Rolling = GetRolling(this.Data[:])
	return this.Rolling
}

func GetStrong(buffer []byte) string {
	return fmt.Sprintf("%x", md5.Sum(buffer))
}

func (this *Chunk) CountStrong() string {
	this.Strong = GetStrong(this.Data[:])
	return this.Strong
}