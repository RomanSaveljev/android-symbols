package signatures

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"reflect"
)

func TestSignaturesAddUnique(t *testing.T) {
	assert := assert.New(t)
	var sig strongSignatures
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
	var sig strongSignatures
	sig = addUnique(sig, []byte("444"))
	sig2 := addUnique(sig, []byte("444"))
	sig3 := addUnique(sig2, []byte("555"))
	assert.True(reflect.DeepEqual(sig, sig2))
	assert.False(reflect.DeepEqual(sig2, sig3))
}

func TestSignaturesGetFromEmpty(t *testing.T) {
	var sigs Signatures
	assert.Nil(t, sigs.Get([]byte{'a'}))
}

func TestSignaturesGetFromNotEmpty(t *testing.T) {
	assert := assert.New(t)
	var sigs Signatures
	sigs.Add([]byte("abc"), []byte("def"))
	existing := sigs.Get([]byte("abc"))
	assert.NotNil(existing)
	assert.True(existing.Has([]byte("def")))
	assert.False(existing.Has([]byte("abc")))
}

func TestSignaturesMultipleStrong(t *testing.T) {
	assert := assert.New(t)
	var signatures Signatures
	signatures.Add([]byte("5"), []byte("zzz"))
	signatures.Add([]byte("5"), []byte("yyy"))
	existing := signatures.Get([]byte("5"))
	assert.NotNil(existing)
	assert.True(existing.Has([]byte("zzz")))
	assert.True(existing.Has([]byte("yyy")))
}
