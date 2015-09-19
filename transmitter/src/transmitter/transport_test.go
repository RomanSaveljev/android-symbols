package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os/exec"
)

func TestProcessTransportCat(t *testing.T) {
	tr, err := NewProcessTransport(exec.Command("cat", "-"))
	assert.Equal(t, nil, err)
	sent := []byte("abc abc abcd")
	n, err := tr.Write(sent)
	assert.Equal(t, nil, err)
	assert.Equal(t, 12, n)
	received := make([]byte, 13)
	n, err = tr.Read(received)
	assert.Equal(t, 12, n)
	assert.Equal(t, sent, received[:n])
	assert.Equal(t, nil, err)
	err = tr.Close()
	assert.Equal(t, nil, err)
}
