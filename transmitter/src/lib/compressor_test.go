package transmitter
/*
import (
    "testing"
    "bytes"
    "github.com/stretchr/testify/assert"
)

func TestCompressorEmptySignatures(t *testing.T) {
	var b bytes.Buffer
	s := NewSignatures()
	compress := NewCompressor(s, 1, &b)
	n, err := compress.Write([]byte{'e'})
	assert.Equal(t, 1, n)
	assert.Equal(t, nil, err)
	err = compress.Close()
	assert.Equal(t, nil, err)
	assert.Equal(t, "AH\n", b.String())
}

func TestCompressorLongerBuffer(t *testing.T) {
	var b bytes.Buffer
	s := NewSignatures()
	compress := NewCompressor(s, 3, &b)
	n, err := compress.Write([]byte("eee"))
	assert.Equal(t, 3, n)
	assert.Equal(t, nil, err)
	err = compress.Close()
	assert.Equal(t, nil, err)
	assert.Equal(t, "AS#E\n", b.String())	
}

func TestCompressorLongerBufferOverflow(t *testing.T) {
	var b bytes.Buffer
	s := NewSignatures()
	compress := NewCompressor(s, 2, &b)
	n, err := compress.Write([]byte("eeezzz390"))
	assert.Equal(t, 9, n)
	assert.Equal(t, nil, err)
	err = compress.Close()
	assert.Equal(t, nil, err)
	assert.Equal(t, "AS#G!H?qA-0E\n", b.String())	
}

// check that signature is applied only when the buffer gets full
func TestCompressorSignatureApplied(t *testing.T) {
	var b bytes.Buffer
	s := NewSignatures()
	// adding-stuff
	s.Add(0xafc2496e, "da86ebed5b7e0b178ac75241c9a72c9d")
	// and-yet-more
	s.Add(0x0857fa33, "56402e1eb6af766c99a77d27f875b3de")
	compress := NewCompressor(s, 12, &b)
	compress.Write([]byte("adding-stuff"))
	compress.Write([]byte("and-yet-more..."))
	compress.Close()
	line, _ := b.ReadString('\n')
	assert.Equal(t, "\tafc2496e/da86ebed5b7e0b178ac75241c9a72c9d\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "\t0857fa33/56402e1eb6af766c99a77d27f875b3de\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "/hSa\n", line)
}

func TestCompressorWeakCollision(t *testing.T) {
	var b bytes.Buffer
	s := NewSignatures()
	// adding-stuff
	s.Add(0xafc2496e, "00000000000000000000000000000000")
	compress := NewCompressor(s, 12, &b)
	compress.Write([]byte("adding-stuff"))
	compress.Close()
	assert.Equal(t, "@:Wn_DJ(PBFEM2-\n", b.String())
}

func TestCompressorSignatureInterspersed(t *testing.T) {
	var b bytes.Buffer
	s := NewSignatures()
	// abc
	s.Add(0x352441c2, "900150983cd24fb0d6963f7d28e17f72")
	// efg
	s.Add(0x512ce803, "7d09898e18511cf7c0c1815d07728d23")
	compress := NewCompressor(s, 3, &b)
	compress.Write([]byte("abcdefg"))
	compress.Close()
	line, _ := b.ReadString('\n')
	assert.Equal(t, "\t352441c2/900150983cd24fb0d6963f7d28e17f72\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "A,\n", line)
	line, _ = b.ReadString('\n')
	assert.Equal(t, "\t512ce803/7d09898e18511cf7c0c1815d07728d23\n", line)
}
*/