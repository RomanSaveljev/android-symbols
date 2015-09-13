package shared

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetFromEmpty(t *testing.T) {
	signatures := NewSignatures()
	assert.Equal(t, len(signatures.Get(3)), 0)
}

func TestGetFromNotEmpty(t *testing.T) {
	signatures := NewSignatures()
	signatures.Add(5, "abc")
	existing := signatures.Get(5)
	assert.Equal(t, len(existing), 1)
	assert.Equal(t, existing[0], "abc")
}

func TestGetSortedSlice(t *testing.T) {
	signatures := NewSignatures()
	signatures.Add(5, "zzz")
	signatures.Add(5, "yyy")
	existing := signatures.Get(5)
	assert.Equal(t, len(existing), 2)
	assert.Equal(t, existing[0], "yyy")
	assert.Equal(t, existing[1], "zzz")
}

func TestGetOneOfFew(t *testing.T) {
	signatures := NewSignatures()
	signatures.Add(387, "111")
	signatures.Add(388, "222")
	signatures.Add(387, "333")
	signatures.Add(388, "444")
	slice387 := signatures.Get(387)
	assert.Equal(t, slice387[0], "111")
	assert.Equal(t, slice387[1], "333")
	slice388 := signatures.Get(388)
	assert.Equal(t, slice388[0], "222")
	assert.Equal(t, slice388[1], "444")
}