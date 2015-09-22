package signatures

import (
	_ "bytes"
	_ "sort"
)

type Signatures interface {
	Add(rolling uint32, strong []byte)
	Get(rolling uint32) StrongSignatures
}

func NewSignatures() Signatures {
	sig := signatures{}
	sig.collection = make(map[uint32]strongSignatures)
	return &sig
}

// Signatures collection is arranged by
type signatures struct {
	//collection []groupedSignatures
	collection map[uint32]strongSignatures
}

func (this *signatures) Add(rolling uint32, strong []byte) {
	val, _ := this.collection[rolling]
	this.collection[rolling] = addUnique(val, strong)
}

func (this *signatures) Get(rolling uint32) StrongSignatures {
	val, _ := this.collection[rolling]
	return val
}
