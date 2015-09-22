package chunk

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"crypto/rand"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
)

/*
func TestCountStrong(t *testing.T) {
	buffer := []byte("123456789abcdef")
	expected := CountStrong(buffer, []byte{})
	for i := 1; i < len(buffer); i++ {
		first, second := buffer[:i], buffer[i:]
		actual := CountStrong(first, second)
		assert.Equal(t, expected, actual)
	}
}
*/

var buf = []byte("sadoiewqckjasnasljdlasudfoiewqurpouqwoiduwqd")
var left = md5.Sum(buf)
var right = md5.Sum(buf)

func BenchmarkEqualBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes.Equal(left[:], right[:])
	}
}

func BenchmarkEqualArrays(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = left == right
	}
}

func BenchmarkSscanf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var result uint32
		fmt.Sscanf("deadbeef", "%x", &result)
		fmt.Sscanf("12345678", "%x", &result)
		fmt.Sscanf("12ab34cd", "%x", &result)
		fmt.Sscanf("fff09ae3", "%x", &result)
	}
}

func BenchmarkRollingFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RollingFromString("deadbeef")
		RollingFromString("12345678")
		RollingFromString("12ab34cd")
		RollingFromString("fff09ae3")
	}
}

func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%08x", 0xdeadbeef)
		fmt.Sprintf("%08x", 0x12345678)
		fmt.Sprintf("%08x", 0x12ab34cd)
		fmt.Sprintf("%08x", 0xfff09ae3)
	}
}

func BenchmarkRollingToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RollingToString(0xdeadbeef)
		RollingToString(0x12345678)
		RollingToString(0x12ab34cd)
		RollingToString(0xfff09ae3)
	}
}

func TestRollingFromString(t *testing.T) {
	assert := assert.New(t)
	ret, err := RollingFromString("deadbeef")
	assert.NoError(err)
	assert.Equal(uint32(0xdeadbeef), ret)
	ret, err = RollingFromString("12345678")
	assert.NoError(err)
	assert.Equal(uint32(0x12345678), ret)
	ret, err = RollingFromString("55559999")
	assert.NoError(err)
	assert.Equal(uint32(0x55559999), ret)
	ret, err = RollingFromString("00000000")
	assert.NoError(err)
	assert.Equal(uint32(0x0), ret)
	ret, err = RollingFromString("123456")
	assert.NoError(err)
	assert.Equal(uint32(0x123456), ret)
	ret, err = RollingFromString("1234")
	assert.NoError(err)
	assert.Equal(uint32(0x1234), ret)
	ret, err = RollingFromString("12")
	assert.NoError(err)
	assert.Equal(uint32(0x12), ret)
}

func TestRollingToString(t *testing.T) {
	assert := assert.New(t)
	ret := RollingToString(uint32(0xdeadbeef))
	assert.Equal("deadbeef", ret)
	ret = RollingToString(uint32(0x12345678))
	assert.Equal("12345678", ret)
	ret = RollingToString(uint32(0x55559999))
	assert.Equal("55559999", ret)
	ret = RollingToString(uint32(0x0))
	assert.Equal("00000000", ret)
	ret = RollingToString(uint32(0x123456))
	assert.Equal("00123456", ret)
	ret = RollingToString(uint32(0x1234))
	assert.Equal("00001234", ret)
	ret = RollingToString(uint32(0x12))
	assert.Equal("00000012", ret)
}

func TestRollerRunningWindow(t *testing.T) {
	assert := assert.New(t)
	roller := Roller{}
	buf := make([]byte, receiver.CHUNK_SIZE) 
	n, _ := rand.Read(buf)
	assert.Equal(receiver.CHUNK_SIZE, n)
	firstFull := CountRolling(buf)
	roller.Calculate(buf)
	assert.Equal(firstFull, roller.Value())
	buf = append(buf[1:], 'a')
	nextFull := CountRolling(buf)
	roller.Next(buf[1], buf[len(buf) - 1])
	assert.Equal(nextFull, roller.Value())
}
