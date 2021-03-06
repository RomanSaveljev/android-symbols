package chunker

import (
	"errors"
	"fmt"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/chunk"
	"github.com/RomanSaveljev/android-symbols/transmitter/mock"
	"github.com/RomanSaveljev/android-symbols/transmitter/signatures"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type matchLength struct {
	length int
}

func (this *matchLength) Matches(x interface{}) bool {
	return len(x.([]byte)) == this.length
}

func (this *matchLength) String() string {
	return fmt.Sprintf("length is %d", this.length)
}

func TestChunkerWrite(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)

	chunker := NewChunker(encoder, receiver)
	err := chunker.Write('a')
	assert.NoError(err)
}

func TestChunkerCloseEmpty(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	encoder.EXPECT().Close()

	chunker := NewChunker(encoder, receiver)
	err := chunker.Close()
	assert.NoError(err)
}

func TestChunkerCloseEncoderError(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	encoder.EXPECT().Close().Return(errors.New("BOO!"))

	chunker := NewChunker(encoder, receiver)
	err := chunker.Close()
	assert.Error(err)
}

func TestChunkerCloseFlushes(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	flushingWrite := encoder.EXPECT().Write([]byte{'a', 'b'}).Times(1).Return(2, nil)
	encoder.EXPECT().Close().After(flushingWrite)

	chunker := NewChunker(encoder, receiver)
	err := chunker.Write('a')
	assert.NoError(err)
	err = chunker.Write('b')
	assert.NoError(err)
	err = chunker.Close()
	assert.NoError(err)
}

func TestChunkerCloseFlushWriteError(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	encoder.EXPECT().Write([]byte{'a'}).Return(0, errors.New("ERROR"))

	chunker := NewChunker(encoder, receiver)
	err := chunker.Write('a')
	assert.NoError(err)
	err = chunker.Close()
	assert.Error(err)
}

func TestChunkerCloseFlushIncompleteWrite(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	errorWrite := encoder.EXPECT().Write([]byte{'a', 'b'}).Return(1, errors.New("ERROR"))
	success := encoder.EXPECT().Write([]byte{'b'}).After(errorWrite).Return(1, nil)
	encoder.EXPECT().Close().After(success)

	chunker := NewChunker(encoder, receiver)
	err := chunker.Write('a')
	assert.NoError(err)
	err = chunker.Write('b')
	assert.NoError(err)
	err = chunker.Close()
	assert.Error(err)
	err = chunker.Close()
	assert.NoError(err)
}

func TestChunkerFlush(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	failure := encoder.EXPECT().Write([]byte{'a', 'b'}).Return(1, errors.New("ERROR"))
	encoder.EXPECT().Write([]byte{'b', 'c'}).After(failure).Return(2, nil)

	chunker := NewChunker(encoder, receiver)
	err := chunker.Write('a')
	assert.NoError(err)
	err = chunker.Write('b')
	assert.NoError(err)
	err = chunker.Flush()
	assert.Error(err)
	err = chunker.Write('c')
	assert.NoError(err)
	err = chunker.Flush()
	assert.NoError(err)
}

func TestChunkerWriteSignature(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	encoder.EXPECT().WriteSignature(uint32(0x123), []byte("abc"))

	chunker := NewChunker(encoder, receiver)
	err := chunker.WriteSignature(0x123, []byte("abc"))
	assert.NoError(err)
}

func TestChunkerFlushBeforeWriteSignature(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	receiver := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	flush := encoder.EXPECT().Write([]byte{'a'}).Return(1, nil)
	encoder.EXPECT().WriteSignature(uint32(0x123), []byte("abc")).After(flush)

	chunker := NewChunker(encoder, receiver)
	err := chunker.Write('a')
	assert.NoError(err)
	err = chunker.WriteSignature(0x123, []byte("abc"))
	assert.NoError(err)
}

func TestChunkerFullBufferCreatesSignature(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sigs := signatures.NewSignatures()

	rcv := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	saveChunk := rcv.EXPECT().SaveChunk(gomock.Any(), gomock.Any(), gomock.Any())
	encoder.EXPECT().WriteSignature(gomock.Any(), gomock.Any()).After(saveChunk)

	chunker := NewChunker(encoder, rcv)
	for i := 0; i < receiver.CHUNK_SIZE; i++ {
		err := chunker.Write(byte(i))
		assert.NoError(err)
	}
}

func TestChunkerBufferEmptiesOnFlush(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rcv := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	encoder.EXPECT().Write(&matchLength{receiver.CHUNK_SIZE/2 + 1}).Times(2).Return(receiver.CHUNK_SIZE/2+1, nil)

	chunker := NewChunker(encoder, rcv)
	for i := 0; i < receiver.CHUNK_SIZE/2+1; i++ {
		err := chunker.Write(byte(i))
		assert.NoError(err)
	}
	err := chunker.Flush()
	assert.NoError(err)
	for i := 0; i < receiver.CHUNK_SIZE/2+1; i++ {
		err := chunker.Write(byte(i))
		assert.NoError(err)
	}
	err = chunker.Flush()
	assert.NoError(err)
}

func TestChunkerDoesNotTransmitExistingChunks(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	buff := make([]byte, receiver.CHUNK_SIZE)
	rolling, strong := chunk.CountRolling(buff), chunk.CountStrong(buff)

	sigs := signatures.NewSignatures()
	sigs.Add(rolling, strong)

	rcv := mock.NewMockReceiver(mockCtrl)
	encoder := mock.NewMockEncoder(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	encoder.EXPECT().WriteSignature(rolling, strong)

	chunker := NewChunker(encoder, rcv)
	for _, b := range buff {
		err := chunker.Write(b)
		assert.NoError(err)
	}
}
