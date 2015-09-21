package transmitter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountStrong(t *testing.T) {
	buffer := []byte("123456789abcdef")
	expected := CountStrong(buffer, []byte{})
	for i := 1; i < len(buffer); i++ {
		first, second := buffer[:i], buffer[i:]
		actual := CountStrong(first, second)
		assert.Equal(t, expected, actual)
	}
}
