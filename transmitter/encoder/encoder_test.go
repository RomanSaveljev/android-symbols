package encoder

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
	err := enc.WriteSignature([]byte{1, 2, 3}, []byte{4, 5, 6, 7, 8})
	assert.NoError(err)
	enc.Close()
	assert.Equal("\t010203/0405060708\n", b.String())
}

func TestEncoderWriteTwoSignatures(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.WriteSignature([]byte{0xa, 0xb}, []byte{0x11, 0x12, 0x13})
	enc.WriteSignature([]byte{0xa, 0xb}, []byte{0x21, 0x22, 0x23})
	enc.Close()
	line, _ := b.ReadString('\n')
	assert.Equal("\t0a0b/111213\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\t0a0b/212223\n", line)
}

func TestEncoderSignatureAndBytes(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Write([]byte{'a'})
	enc.WriteSignature([]byte{0xff}, []byte{0x8c, 0x17, 0xa6, 0x83})
	enc.Write([]byte{'b'})
	enc.Write([]byte{'c'})
	enc.WriteSignature([]byte{0xff, 0xaa}, []byte{0xa7, 0x0e, 0x2e, 0x20, 0x8f})
	enc.Write([]byte{'d'})
	enc.Write([]byte{'e'})
	enc.Write([]byte{'f'})
	enc.WriteSignature([]byte{0xff, 0xaa, 0xbb}, []byte{0xe8, 0x45, 0x03, 0x41, 0xf1})
	enc.Close()
	line, _ := b.ReadString('\n')
	assert.Equal("@/\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tff/8c17a683\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("@Uf\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tffaa/a70e2e208f\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("A7]?\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal("\tffaabb/e8450341f1\n", line)
	line = b.String()
	assert.Equal("", line)
}
