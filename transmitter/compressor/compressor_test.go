package compressor

import (
	"crypto/rand"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/mock"
	"github.com/RomanSaveljev/android-symbols/transmitter/signatures"
	"github.com/RomanSaveljev/android-symbols/transmitter/chunk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/hex"
)

func TestCompressorWriteStuffsTheBuffer(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	buff := make([]byte, receiver.CHUNK_SIZE/2-1)
	assert.Equal(receiver.CHUNK_SIZE/2-1, len(buff))

	rcv := mock.NewMockReceiver(mockCtrl)
	chunker := mock.NewMockChunker(mockCtrl)

	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write(buff)
	assert.Equal(receiver.CHUNK_SIZE/2-1, n)
	assert.NoError(err)
	n, err = compressor.Write(buff)
	assert.Equal(receiver.CHUNK_SIZE/2-1, n)
	assert.NoError(err)
	n, err = compressor.Write([]byte{'a'})
	assert.Equal(1, n)
	assert.NoError(err)
}

func randomChunk(t *testing.T) (c receiver.Chunk, rolling uint32, strong []byte) {
	n, err := rand.Read(c.Data[:])
	assert.NoError(t, err)
	assert.Equal(t, len(c.Data), n)
	rolling = chunk.CountRolling(c.Data[:])
	c.Rolling = chunk.RollingToString(rolling)
	strong = chunk.CountStrong(c.Data[:])
	c.Strong = hex.EncodeToString(strong)
	return
}

func TestCompressorFullBufferWritesSignature(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	chunk, rolling, strong := randomChunk(t)
	sigs := signatures.NewSignatures()
	sigs.Add(rolling, strong)

	rcv := mock.NewMockReceiver(mockCtrl)
	chunker := mock.NewMockChunker(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	chunker.EXPECT().WriteSignature(rolling, strong).Times(1)

	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write(chunk.Data[:])
	assert.NoError(err)
	assert.Equal(len(chunk.Data), n)
	// compressor's buffer is empty now
	n, err = compressor.Write(chunk.Data[:len(chunk.Data)-1])
	assert.NoError(err)
	assert.Equal(len(chunk.Data)-1, n)
}

func TestCompressorRecognizeSignatureAtOffset(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	chunk, rolling, strong := randomChunk(t)
	sigs := signatures.NewSignatures()
	sigs.Add(rolling, strong)

	rcv := mock.NewMockReceiver(mockCtrl)
	chunker := mock.NewMockChunker(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	first := chunker.EXPECT().Write(byte('a'))
	chunker.EXPECT().WriteSignature(rolling, strong).After(first)

	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write([]byte{'a'})
	assert.NoError(err)
	assert.Equal(1, n)
	n, err = compressor.Write(chunk.Data[:])
	assert.NoError(err)
	assert.Equal(len(chunk.Data), n)
}

func TestCompressorOverlappingSignatures(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c, rolling, strong := randomChunk(t)
	sigs := signatures.NewSignatures()
	sigs.Add(rolling, strong)
	var overlap receiver.Chunk
	copy(overlap.Data[:], append(c.Data[1:], c.Data[0]+1))
	overlapRolling := chunk.CountRolling(overlap.Data[:])
	overlap.Rolling = chunk.RollingToString(overlapRolling)
	overlapStrong := chunk.CountStrong(overlap.Data[:])
	overlap.Strong = hex.EncodeToString(overlapStrong)
	sigs.Add(overlapRolling, overlapStrong)

	rcv := mock.NewMockReceiver(mockCtrl)
	chunker := mock.NewMockChunker(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	gomock.InOrder(
		chunker.EXPECT().WriteSignature(rolling, strong),
		chunker.EXPECT().Write(overlap.Data[len(overlap.Data)-1]),
		chunker.EXPECT().Close(),
	)

	buffer := append(c.Data[:1], overlap.Data[:]...)
	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write(buffer)
	assert.NoError(err)
	assert.Equal(len(buffer), n)
	err = compressor.Close()
	assert.NoError(err)
}

func TestCompressorNoMatchingSignature(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c, rolling, strong := randomChunk(t)
	sigs := signatures.NewSignatures()
	sigs.Add(rolling, append(strong, 'b'))

	rcv := mock.NewMockReceiver(mockCtrl)
	chunker := mock.NewMockChunker(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	write := chunker.EXPECT().Write(c.Data[0])

	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write(c.Data[:])
	assert.NoError(err)
	assert.Equal(len(c.Data), n)

	write = chunker.EXPECT().Write(gomock.Any()).After(write).Times(len(c.Data) - 1)
	chunker.EXPECT().Close().After(write)
	err = compressor.Close()
	assert.NoError(err)
}

func TestCompressorCloseWritesRemainingBytes(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rcv := mock.NewMockReceiver(mockCtrl)
	chunker := mock.NewMockChunker(mockCtrl)

	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write([]byte{'a', 'b', 'c'})
	assert.NoError(err)
	assert.Equal(3, n)
	gomock.InOrder(
		chunker.EXPECT().Write(byte('a')),
		chunker.EXPECT().Write(byte('b')),
		chunker.EXPECT().Write(byte('c')),
		chunker.EXPECT().Close(),
	)
	err = compressor.Close()
	assert.NoError(err)
}
