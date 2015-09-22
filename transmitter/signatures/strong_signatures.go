package signatures

import (
	"bytes"
)

type StrongSignatures interface {
	Has(strong []byte) bool
	IsEmpty() bool
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

func (this strongSignatures) IsEmpty() bool {
	return len(this) == 0
}

func addUnique(this strongSignatures, strong []byte) strongSignatures {
	if !this.Has(strong) {
		this = append(this, append([]byte{}, strong...))
	}
	return this
}
