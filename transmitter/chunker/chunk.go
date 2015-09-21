package chunker

import (
	//"github.com/Redundancy/go-sync/rollsum"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	//"hash/crc32"
	"crypto/md5"
	//"fmt"
	"encoding/binary"
)

type Roller struct {
	a, b uint32
	first byte
	calculatedFirst bool
	result [4]byte
}

func (this *Roller) Value() []byte {
	value := uint32(this.a) + (uint32(this.b) << 16)
	b := this.result[:]
	binary.LittleEndian.PutUint32(b, value)
	return b
}

func (this *Roller) Next(first, last byte) {
	if !this.calculatedFirst {
		panic("Must do Calculate() first")
	}
	this.b = this.b - (receiver.CHUNK_SIZE + 1) * uint32(first) + this.a
	this.a = this.a - uint32(this.first) + uint32(last)
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
		this.b += uint32(receiver.CHUNK_SIZE - i + 1) * uint32(b)
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

func CountRolling(buffer[] byte) []byte {
	roller := Roller{}
	roller.Calculate(buffer)
	return roller.Value()
}

func CountStrong(buffer []byte) []byte {
	result := md5.Sum(buffer)
	return result[:]
}
