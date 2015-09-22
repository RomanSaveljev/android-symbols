package signatures

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSignaturesAddUnique(t *testing.T) {
	assert := assert.New(t)
	var sig StrongSignatures
	sig = addUnique(sig, []byte("444"))
	sig = addUnique(sig, []byte("111"))
	sig = addUnique(sig, []byte("222"))
	sig = addUnique(sig, []byte("555"))
	assert.True(sig.Has([]byte("111")))
	assert.True(sig.Has([]byte("222")))
	assert.True(sig.Has([]byte("444")))
	assert.True(sig.Has([]byte("555")))
	assert.False(sig.Has([]byte("aaa")))
}

func TestSignaturesAddUniqueNoDuplicates(t *testing.T) {
	assert := assert.New(t)
	var sig StrongSignatures
	sig = addUnique(sig, []byte("444"))
	sig2 := addUnique(sig, []byte("444"))
	sig3 := addUnique(sig2, []byte("555"))
	assert.True(reflect.DeepEqual(sig, sig2))
	assert.False(reflect.DeepEqual(sig2, sig3))
}

func TestSignaturesGetFromEmpty(t *testing.T) {
	sigs := NewSignatures()
	assert.Empty(t, sigs.Get(0xa))
}

func TestSignaturesGetFromNotEmpty(t *testing.T) {
	assert := assert.New(t)
	sigs := NewSignatures()
	sigs.Add(0xabc, []byte("def"))
	existing := sigs.Get(0xabc)
	assert.NotEmpty(existing)
	assert.True(existing.Has([]byte("def")))
	assert.False(existing.Has([]byte("abc")))
}

func TestSignaturesMultipleStrong(t *testing.T) {
	assert := assert.New(t)
	signatures := NewSignatures()
	signatures.Add(0x5, []byte("zzz"))
	signatures.Add(0x5, []byte("yyy"))
	existing := signatures.Get(0x5)
	assert.NotEmpty(existing)
	assert.True(existing.Has([]byte("zzz")))
	assert.True(existing.Has([]byte("yyy")))
}

func TestSignaturesMultipleRolling(t *testing.T) {
	assert := assert.New(t)
	signatures := NewSignatures()
	signatures.Add(0xabcd, []byte("def0"))
	group := signatures.Get(0xabcd)
	assert.NotNil(group)
	signatures.Add(0x1234, []byte("4567"))
	group = signatures.Get(0x1234)
	assert.NotEmpty(group)
}

func TestSignaturesSearchRolling(t *testing.T) {
	assert := assert.New(t)
	signatures := NewSignatures()
	signatures.Add(0xdd176afb, []byte{1})
	signatures.Add(0xdd4582ea, []byte{2})
	signatures.Add(0xdd5eb264, []byte{3})
	signatures.Add(0xdd72e6b2, []byte{4})
	signatures.Add(0xdd78f790, []byte{5})
	signatures.Add(0xdd897335, []byte{6})
	candidates := signatures.Get(0xdd4c51d8)
	assert.Empty(candidates)
}
