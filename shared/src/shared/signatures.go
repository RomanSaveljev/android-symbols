package shared

import (
	"sort"
)

// Signatures collection is arranged by
type Signatures struct {
	collection map[uint32][]string
}

func NewSignatures() *Signatures {
	signatures := Signatures{make(map[uint32][]string)}
	return &signatures
}

func (this *Signatures) Add(rolling uint32, sig string) {
	_, exists := this.collection[rolling]
	if !exists {
		existing := make([]string, 0, 1)
		this.collection[rolling] = existing
	}
	if sort.SearchStrings(this.collection[rolling], sig) != -1 {
		this.collection[rolling] = append(this.collection[rolling], sig)
		sort.Strings(this.collection[rolling])
	}
}

func (this *Signatures) Get(rolling uint32) []string {
	existing, exists := this.collection[rolling]
	if !exists {
		existing = []string{}
	}
	return existing
}
