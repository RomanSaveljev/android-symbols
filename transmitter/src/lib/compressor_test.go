package transmitter

import (
	"crypto/rand"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/mock"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/signatures"
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

	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	chunker := mock_transmitter.NewMockChunker(mockCtrl)

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

func randomChunk(t *testing.T) (chunk receiver.Chunk, rolling, strong []byte) {
	n, err := rand.Read(chunk.Data[:])
	assert.NoError(t, err)
	assert.Equal(t, len(chunk.Data), n)
	rolling = CountRolling(chunk.Data[:], []byte{})
	chunk.Rolling = hex.EncodeToString(rolling)
	strong = CountStrong(chunk.Data[:], []byte{})
	chunk.Strong = hex.EncodeToString(strong)
	return
}

func TestCompressorFullBufferWritesSignature(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	chunk, rolling, strong := randomChunk(t)
	sigs := new(signatures.Signatures)
	sigs.Add(rolling, strong)

	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	chunker := mock_transmitter.NewMockChunker(mockCtrl)
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
	sigs := new(signatures.Signatures)
	sigs.Add(rolling, strong)

	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	chunker := mock_transmitter.NewMockChunker(mockCtrl)
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

	chunk, rolling, strong := randomChunk(t)
	sigs := new(signatures.Signatures)
	sigs.Add(rolling, strong)
	var overlap receiver.Chunk
	copy(overlap.Data[:], append(chunk.Data[1:], chunk.Data[0]+1))
	overlapRolling := CountRolling(overlap.Data[:], []byte{})
	overlap.Rolling = hex.EncodeToString(overlapRolling)
	overlapStrong := CountStrong(overlap.Data[:], []byte{})
	overlap.Strong = hex.EncodeToString(overlapStrong)
	sigs.Add(overlapRolling, overlapStrong)

	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	chunker := mock_transmitter.NewMockChunker(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	gomock.InOrder(
		chunker.EXPECT().WriteSignature(rolling, strong),
		chunker.EXPECT().Write(overlap.Data[len(overlap.Data)-1]),
		chunker.EXPECT().Close(),
	)

	buffer := append(chunk.Data[:1], overlap.Data[:]...)
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

	chunk, rolling, strong := randomChunk(t)
	sigs := new(signatures.Signatures)
	sigs.Add(rolling, append(strong, 'b'))

	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	chunker := mock_transmitter.NewMockChunker(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	write := chunker.EXPECT().Write(chunk.Data[0])

	compressor := NewCompressor(chunker, rcv)
	n, err := compressor.Write(chunk.Data[:])
	assert.NoError(err)
	assert.Equal(len(chunk.Data), n)

	write = chunker.EXPECT().Write(gomock.Any()).After(write).Times(len(chunk.Data) - 1)
	chunker.EXPECT().Close().After(write)
	err = compressor.Close()
	assert.NoError(err)
}

func TestCompressorCloseWritesRemainingBytes(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	chunker := mock_transmitter.NewMockChunker(mockCtrl)

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
