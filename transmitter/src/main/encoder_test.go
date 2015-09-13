package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"bytes"
)

func TestEncoderWriteByte(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	err := enc.Write(1)
	assert.Equal(t, nil, err)
}

func TestEncoderWriteOne(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	enc.Write('a')
	enc.Close()
	assert.Equal(t, "@/\n", b.String())
}

func TestEncoderWriteTwo(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	enc.Write('a')
	enc.Write('b')
	enc.Close()
	assert.Equal(t, "@:B\n", b.String())
}

func TestEncoderWriteThree(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	enc.Write('a')
	enc.Write('b')
	enc.Write('c')
	enc.Close()
	assert.Equal(t, "@:E^\n", b.String())
}

func TestEncoderWriteBytes(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	enc.Write('a')
	enc.Write('b')
	enc.Write('c')
	enc.Write('d')
	err := enc.Close()
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("@:E_W\n"), b.Bytes())
}

func TestEncoderWriteSignature(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	err := enc.WriteSignature(0xdeadbeef, "abcdef012345678")
	assert.Equal(t, nil, err)
	enc.Close()
	assert.Equal(t, "\tdeadbeef/abcdef012345678\n", string(b.Bytes()))
}

func TestEncoderWriteTwoSignatures(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	enc.WriteSignature(0xdeadbeef, "abcdef012345678")
	enc.WriteSignature(0xdeadbeef, "abcdef012345679")
	enc.Close()
	line, _ := b.ReadString('\n')
	assert.Equal(t, "\tdeadbeef/abcdef012345678\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "\tdeadbeef/abcdef012345679\n", line)
}

func TestEncoderSignatureAndBytes(t *testing.T) {
	var b bytes.Buffer
	enc := newEncoder(&b)
	enc.Write('a')
	enc.WriteSignature(0xff, "8c17a6833de2c1766302dd7477ee4a20")
	enc.Write('b')
	enc.Write('c')
	enc.WriteSignature(0xffaa, "a70e2e208f3c5c9d7cd52e148f40178b")
	enc.Write('d')
	enc.Write('e')
	enc.Write('f')
	enc.WriteSignature(0xffaabb, "e8450341f161f65372fbd784fe28c8f5")
	enc.Close()
	line, _ := b.ReadString('\n')
	assert.Equal(t, "@/\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "\t000000ff/8c17a6833de2c1766302dd7477ee4a20\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "@Uf\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "\t0000ffaa/a70e2e208f3c5c9d7cd52e148f40178b\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "A7]?\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "\t00ffaabb/e8450341f161f65372fbd784fe28c8f5\n", line)
	line = b.String()
	assert.Equal(t, "", line)
}