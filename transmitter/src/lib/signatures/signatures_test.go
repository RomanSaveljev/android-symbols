package signatures

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"encoding/hex"
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

func TestSignaturesMultipleRolling(t *testing.T) {
	assert := assert.New(t)
	var signatures Signatures
	signatures.Add([]byte("abcd"), []byte("def0"))
	group := signatures.Get([]byte("abcd"))
	assert.NotNil(group)
	signatures.Add([]byte("1234"), []byte("4567"))
	group = signatures.Get([]byte("1234"))
	assert.NotNil(group)
}

func TestSignaturesSearchRolling(t *testing.T) {
	decode := func(s string) []byte {
		ret, _ := hex.DecodeString(s)
		return ret
	}
	assert := assert.New(t)
	var signatures Signatures
	signatures.Add(decode("dd176afb888331c1b08efb324200f277"), []byte{1})
	signatures.Add(decode("dd4582ea82e6e153ddfc22710bea7fef"), []byte{2})
	signatures.Add(decode("dd5eb2641db30a4d24934c3fc873aeba"), []byte{3})
	signatures.Add(decode("dd72e6b2dca00d9ce5734b3ec90c96a2"), []byte{4})
	signatures.Add(decode("dd78f7905d03d8af54feeea290a7ce80"), []byte{5})
	signatures.Add(decode("dd897335619c9d02f9e3a82fcc48bb4d"), []byte{6})
	candidates := signatures.Get(decode("dd4c51d8c594340bca4898155feaed81"))
	assert.Nil(candidates)
}