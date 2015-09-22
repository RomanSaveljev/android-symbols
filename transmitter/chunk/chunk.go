package chunk

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
)

type Roller struct {
	a, b            uint32
	first           byte
	calculatedFirst bool
	result          [4]byte
}

func (this *Roller) Value() uint32 {
	return uint32(this.a) + (uint32(this.b) << 16)
}

func (this *Roller) Next(first, last byte) {
	if !this.calculatedFirst {
		panic("Must do Calculate() first")
	}	
	this.a = this.a - uint32(this.first) + uint32(last)
	this.b = this.b - receiver.CHUNK_SIZE * uint32(this.first) + this.a
	this.normalize()
	this.first = first
}

func (this *Roller) Calculate(buffer []byte) {
	// this is from tech_report.tex distributed along the rsync source code
	if len(buffer) != receiver.CHUNK_SIZE {
		panic("Roller only works with predefined block size")
	}
	this.a = 0
	this.b = 0
	for i, b := range buffer {
		this.a += uint32(b)
		this.b += uint32(receiver.CHUNK_SIZE - i) * uint32(b)
	}
	this.normalize()
	this.calculatedFirst = true
	this.first = buffer[0]
}

func (this *Roller) Calculated() bool {
	return this.calculatedFirst
}

func (this *Roller) normalize() {
	this.a &= 0xffff
	this.b &= 0xffff
}

func CountRolling(buffer []byte) uint32 {
	roller := Roller{}
	roller.Calculate(buffer)
	return roller.Value()
}

func CountStrong(buffer []byte) []byte {
	result := md5.Sum(buffer)
	// based on benchmarking bytes.Equal(a, b) is the same as a == b
	// so, it is ok to pass around strong signature as a slice
	return result[:]
}

func RollingFromString(s string) (ret uint32, err error) {
	var buf []byte
	if buf, err = hex.DecodeString(s); err == nil {
		for len(buf) < 4 {
			buf = append([]byte{0}, buf...)
		}
		ret = binary.BigEndian.Uint32(buf)
	}
	return
}

func RollingToString(r uint32) string {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, r)
	return hex.EncodeToString(buf)
}
