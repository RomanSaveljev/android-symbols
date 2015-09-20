package transmitter

import (
	//"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	//"hash/crc32"
	"crypto/md5"
	//"fmt"
)

func CountRolling(buffer []byte) []byte {
	return buffer[:16]
}

func CountStrong(buffer []byte) []byte {
	ret := md5.Sum(buffer)
	return ret[:]
}
