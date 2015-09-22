package receiver

import (
	"crypto/rand"
	_ "encoding/binary"
	"encoding/hex"
	_ "fmt"
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReceiver(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock.NewMockClient(mockCtrl)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "file-token"),
		client.EXPECT().Call("file-token.StartStream", gomock.Not(nil), gomock.Any()).SetArg(2, "stream-token"),
		client.EXPECT().Call("stream-token.Write", []byte("abc"), gomock.Any()).SetArg(2, 3),
		client.EXPECT().Call("stream-token.Close", gomock.Not(nil), gomock.Any()),
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
	chunk.Rolling = "01020304"
	chunk.Strong = hex.EncodeToString([]byte{3, 4})
	n, err := rand.Read(chunk.Data[:])
	assert.NoError(err)
	assert.Equal(len(chunk.Data), n)

	client := mock.NewMockClient(mockCtrl)
	client.EXPECT().Call("tkn.NextSignature", gomock.Any(), gomock.Any()).Return(io.EOF)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "tkn"),
		client.EXPECT().Call("tkn.SaveChunk", chunk, gomock.Any()),
	)

	rcv, err := NewReceiver("/a/b/c/d.txt", client)
	assert.NoError(err)
	err = rcv.SaveChunk(0x01020304, []byte{3, 4}, chunk.Data[:])
	assert.NoError(err)
	sigs, err := rcv.Signatures()
	assert.NoError(err)
	assert.True(sigs.Get(0x01020304).Has([]byte{3, 4}))
	err = rcv.Close()
	assert.NoError(err)
}

func TestReceiverCollectSignatures(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	one := receiver.Signature{Rolling: "abcd", Strong: "def0"}
	two := receiver.Signature{Rolling: "1234", Strong: "4567"}
	three := receiver.Signature{Rolling: "1234", Strong: "789a"}

	client := mock.NewMockClient(mockCtrl)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "tkn"),
		client.EXPECT().Call("tkn.NextSignature", gomock.Not(nil), gomock.Any()).SetArg(2, one),
		client.EXPECT().Call("tkn.NextSignature", gomock.Not(nil), gomock.Any()).SetArg(2, two),
		client.EXPECT().Call("tkn.NextSignature", gomock.Not(nil), gomock.Any()).SetArg(2, three),
		client.EXPECT().Call("tkn.NextSignature", gomock.Not(nil), gomock.Any()).Return(io.EOF),
	)

	rcv, err := NewReceiver("/a/b/c/d.txt", client)
	assert.NoError(err)
	sigs, err := rcv.Signatures()
	assert.NoError(err)
	assert.True(sigs.Get(0xabcd).Has([]byte{0xde, 0xf0}))
	assert.True(sigs.Get(0x1234).Has([]byte{0x45, 0x67}))
	assert.True(sigs.Get(0x1234).Has([]byte{0x78, 0x9a}))
	err = rcv.Close()
	assert.NoError(err)
}

func TestReceiverGetCachedSignatures(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mock.NewMockClient(mockCtrl)
	gomock.InOrder(
		client.EXPECT().Call("Synchronizer.StartFile", "/a/b/c/d.txt", gomock.Any()).SetArg(2, "tkn"),
		client.EXPECT().Call("tkn.NextSignature", gomock.Not(nil), gomock.Any()).Return(io.EOF),
	)

	rcv, err := NewReceiver("/a/b/c/d.txt", client)
	assert.NoError(err)
	sigs, err := rcv.Signatures()
	assert.NoError(err)
	candidates := sigs.Get(0x1234)
	assert.Nil(candidates)
	sigs, err = rcv.Signatures()
	assert.NoError(err)
	candidates = sigs.Get(0x4567)
	assert.Nil(candidates)
}
