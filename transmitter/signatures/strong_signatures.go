package signatures

import (
	"bytes"
	"sort"
)

type StrongSignatures interface {
	Has(strong []byte) bool
}

type strongSignatures [][]byte

func (this strongSignatures) Has(strong []byte) bool {
	searchStrong := func(i int) bool {
		return bytes.Compare(this[i], strong) >= 0
	}
	n := len(this)
	idx := sort.Search(n, searchStrong)
	ret := idx != n && bytes.Equal(strong, this[idx])
	return ret
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
		sort.Sort(this)
	}
	return this
}
