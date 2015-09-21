package signatures

import (
	"bytes"
)

type StrongSignatures interface {
	Has(strong []byte) bool
}

type strongSignatures [][]byte

func (this strongSignatures) Has(strong []byte) bool {
	for _, s := range this {
		if bytes.Equal(strong, s) {
			return true
		}
	}
	return false
}

func (this strongSignatures) Len() int {
	return len(this)
}

func (this strongSignatures) Less(i, j int) bool {
	return bytes.Compare(this[i], this[j]) < 0
}

func (this strongSignatures) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func addUnique(this strongSignatures, strong []byte) strongSignatures {
	if !this.Has(strong) {
		this = append(this, append([]byte{}, strong...))
	}
	return this
}
