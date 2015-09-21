package transmitter

import (
	//"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	//"hash/crc32"
	"crypto/md5"
	//"fmt"
)

func CountRolling(buffer, extra []byte) []byte {
	// TODO: this is a temporary implementation and we must walk through
	// a complete buffer
	if len(buffer) < 16 {
		ret := make([]byte, 16)
		copy(ret, buffer)
		copy(ret[len(buffer):], extra)
		return ret
	} else {
		return buffer[:16]
	}
}

func CountStrong(buffer, extra []byte) []byte {
	hash := md5.New()
	hash.Write(buffer)
	hash.Write(extra)
	return hash.Sum(nil)
}
