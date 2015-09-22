package chunk

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
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
}
