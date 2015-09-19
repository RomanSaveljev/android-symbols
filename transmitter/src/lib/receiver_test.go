package transmitter

import (
	"crypto/rand"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReceiver(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock_transmitter.NewMockClient(mockCtrl)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "file-token"),
		client.EXPECT().Call("file-token.StartStream", gomock.Any(), gomock.Any()).SetArg(2, "stream-token"),
		client.EXPECT().Call("stream-token.Write", []byte("abc"), gomock.Any()).SetArg(2, 3),
		client.EXPECT().Call("stream-token.Close", gomock.Any(), gomock.Any()),
	)

	rcv, err := NewReceiver("/a/b/c/d.txt", client)
	assert.NoError(err)
	n, err := rcv.Write([]byte("abc"))
	assert.NoError(err)
	assert.Equal(3, n)
	err = rcv.Close()
	assert.NoError(err)
}

func TestReceiverSaveChunk(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var chunk receiver.Chunk
	chunk.Rolling = "abc"
	chunk.Strong = "def"
	n, err := rand.Read(chunk.Data[:])
	assert.NoError(err)
	assert.Equal(len(chunk.Data), n)

	client := mock_transmitter.NewMockClient(mockCtrl)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "tkn"),
		client.EXPECT().Call("tkn.SaveChunk", chunk, gomock.Any()),
	)

	rcv, err := NewReceiver("/a/b/c/d.txt", client)
	assert.NoError(err)
	err = rcv.SaveChunk(&chunk)
	assert.NoError(err)
	err = rcv.Close()
	assert.NoError(err)
}

func TestReceiverCollectSignatures(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	one := receiver.Signature{Rolling: "abc", Strong: "def"}
	two := receiver.Signature{Rolling: "123", Strong: "456"}
	three := receiver.Signature{Rolling: "123", Strong: "789"}

	client := mock_transmitter.NewMockClient(mockCtrl)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "tkn"),
		client.EXPECT().Call("tkn.NextSignature", gomock.Any(), gomock.Any()).SetArg(2, one),
		client.EXPECT().Call("tkn.NextSignature", gomock.Any(), gomock.Any()).SetArg(2, two),
		client.EXPECT().Call("tkn.NextSignature", gomock.Any(), gomock.Any()).SetArg(2, three),
		client.EXPECT().Call("tkn.NextSignature", gomock.Any(), gomock.Any()).Return(io.EOF),
	)

	rcv, err := NewReceiver("/a/b/c/d.txt", client)
	assert.NoError(err)
	sigs, err := rcv.Signatures()
	assert.NoError(err)
	assert.Equal("def", sigs.Get("abc")[0])
	assert.Equal("456", sigs.Get("123")[0])
	assert.Equal("789", sigs.Get("123")[1])
	err = rcv.Close()
	assert.NoError(err)
}
