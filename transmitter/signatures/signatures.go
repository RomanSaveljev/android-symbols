package signatures

import (
	"bytes"
	"sort"
)

type groupedSignatures struct {
	rolling []byte
	strongSignatures  strongSignatures
}

type byRollingSorter struct {
	collection []groupedSignatures
}

func (this byRollingSorter) Len() int {
	return len(this.collection)
}

func (this byRollingSorter) Less(i, j int) bool {
	return bytes.Compare(this.collection[i].rolling, this.collection[j].rolling) < 0
}

func (this byRollingSorter) Swap(i, j int) {
	this.collection[i], this.collection[j] = this.collection[j], this.collection[i]
}

// Signatures collection is arranged by
type Signatures struct {
	collection []groupedSignatures
}

func (this *Signatures) findByRolling(rolling []byte) *groupedSignatures {
	searchRolling := func(i int) bool {
		return bytes.Compare(this.collection[i].rolling, rolling) >= 0
	}
	idxRolling := sort.Search(len(this.collection), searchRolling)
	if idxRolling != len(this.collection) && bytes.Equal(rolling, this.collection[idxRolling].rolling) {
		return &this.collection[idxRolling]
	} else {
		return nil
	}
}

func (this *Signatures) Add(rolling, strong []byte) {
	group := this.findByRolling(rolling)
	if group == nil {
		group = new(groupedSignatures)
		group.rolling = append([]byte{}, rolling...)
		group.strongSignatures = addUnique(group.strongSignatures, strong)
		this.collection = append(this.collection, *group)
		sort.Sort(byRollingSorter{this.collection})
	} else {
		group.strongSignatures = addUnique(group.strongSignatures, strong)
	}
}

func (this *Signatures) Get(rolling []byte) (ret StrongSignatures) {
	sig := this.findByRolling(rolling)
	if sig != nil {
		ret = &sig.strongSignatures
	}
	return
}
