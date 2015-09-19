package transmitter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"bytes"
)

func TestEncoderWriteByte(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	n, err := enc.Write([]byte{1})
	assert.NoError(err)
	assert.Equal(1, n)
}

func TestEncoderWriteOne(t *testing.T) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Write([]byte{'a'})
	enc.Close()
	assert.Equal(t, "@/\n", b.String())
}

func TestEncoderWriteTwo(t *testing.T) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Write([]byte{'a', 'b'})
	enc.Close()
	assert.Equal(t, "@:B\n", b.String())
}

func TestEncoderWriteThree(t *testing.T) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Write([]byte{'a', 'b'})
	enc.Write([]byte{'c'})
	enc.Close()
	assert.Equal(t, "@:E^\n", b.String())
}

func TestEncoderWriteBytes(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Write([]byte{'a', 'b', 'c', 'd'})
	err := enc.Close()
	assert.NoError(err)
	assert.Equal("@:E_W\n", b.String())
}

func TestEncoderWriteSignature(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	err := enc.WriteSignature("deadbeef", "abcdef012345678")
	assert.NoError(err)
	enc.Close()
	assert.Equal("\tdeadbeef/abcdef012345678\n", b.String())
}

func TestEncoderWriteTwoSignatures(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.WriteSignature("deadbeef", "abcdef012345678")
	enc.WriteSignature("deadbeef", "abcdef012345679")
	enc.Close()
	line, _ := b.ReadString('\n')
	assert.Equal("\tdeadbeef/abcdef012345678\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tdeadbeef/abcdef012345679\n", line)
}

func TestEncoderSignatureAndBytes(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Write([]byte{'a'})
	enc.WriteSignature("ff", "8c17a6833de2c1766302dd7477ee4a20")
	enc.Write([]byte{'b'})
	enc.Write([]byte{'c'})
	enc.WriteSignature("ffaa", "a70e2e208f3c5c9d7cd52e148f40178b")
	enc.Write([]byte{'d'})
	enc.Write([]byte{'e'})
	enc.Write([]byte{'f'})
	enc.WriteSignature("ffaabb", "e8450341f161f65372fbd784fe28c8f5")
	enc.Close()
	line, _ := b.ReadString('\n')
	assert.Equal("@/\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tff/8c17a6833de2c1766302dd7477ee4a20\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("@Uf\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tffaa/a70e2e208f3c5c9d7cd52e148f40178b\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("A7]?\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tffaabb/e8450341f161f65372fbd784fe28c8f5\n", line)
	line = b.String()
	assert.Equal("", line)
}
