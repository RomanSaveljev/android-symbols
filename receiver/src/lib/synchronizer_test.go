package receiver

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/rpc"
	"path"
	"testing"
	//"log"
	"fmt"
	"os"
)

type pipePair struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	io.WriteCloser
}

func (this *pipePair) Read(p []byte) (int, error) {
	return this.reader.Read(p)
}

func (this *pipePair) Write(p []byte) (int, error) {
	return this.writer.Write(p)
}

func (this *pipePair) Close() error {
	this.writer.Close()
	return this.reader.Close()
}

type context struct {
	assert     *assert.Assertions
	in         pipePair
	out        pipePair
	client     *rpc.Client
	baseFolder string
}

func (this *context) createBaseFolder() {
	var err error
	this.baseFolder, err = ioutil.TempDir("", "test_")
	this.assert.NoError(err)
}

func (this *context) createServer() {
	rpc.DefaultServer = rpc.NewServer()
	this.in.reader, this.out.writer = io.Pipe()
	this.out.reader, this.in.writer = io.Pipe()
	go RunSynchronizerService(&this.out)
	this.client = rpc.NewClient(&this.in)
}

func (this *context) startFile(name string) (pathname string, token string) {
	pathname = path.Join(this.baseFolder, name)
	err := this.client.Call("Synchronizer.StartFile", pathname, &token)
	this.assert.NoError(err)
	return
}

func (this *context) nextSignature(token string) Signature {
	var signature Signature
	err := this.client.Call(fmt.Sprint(token, ".NextSignature"), 0, &signature)
	this.assert.NoError(err)
	return signature
}

func (this *context) startStream(token string) string {
	var stream string
	err := this.client.Call(fmt.Sprint(token, ".StartStream"), 0, &stream)
	this.assert.NoError(err)
	return stream
}

func (this *context) writeStream(token string, data []byte) {
	var n int
	err := this.client.Call(fmt.Sprint(token, ".Write"), data, &n)
	this.assert.NoError(err)
	this.assert.Equal(len(data), n)
}

func (this *context) closeStream(token string) {
	err := this.client.Call(fmt.Sprint(token, ".Close"), 0, nil)
	this.assert.NoError(err)
}

func TestSynchronizerStartFile(t *testing.T) {
	assert := assert.New(t)
	ctx := context{assert: assert}
	ctx.createBaseFolder()
	ctx.createServer()
	defer ctx.client.Close()
	pathname, token := ctx.startFile("test")
	_, err := os.Stat(pathname)
	assert.False(os.IsNotExist(err))
	assert.NotEqual(0, len(token))
}

func TestSynchronizerFileNextSignature(t *testing.T) {
	assert := assert.New(t)
	ctx := context{assert: assert}
	ctx.createBaseFolder()

	sigPath := path.Join(ctx.baseFolder, "test", "ffffeeee")
	err := os.MkdirAll(sigPath, os.ModeDir|os.ModePerm)
	assert.NoError(err)
	err = ioutil.WriteFile(path.Join(sigPath, "6a900ff954cc07f2ea7b0fa2b862c480"), []byte("123"), os.ModePerm)
	assert.NoError(err)

	ctx.createServer()
	defer ctx.client.Close()
	_, token := ctx.startFile("test")
	signature := ctx.nextSignature(token)
	assert.Equal("ffffeeee", signature.Rolling)
	assert.Equal("6a900ff954cc07f2ea7b0fa2b862c480", signature.Strong)
	ctx.client.Close()
}

func TestSynchronizerFileStartStream(t *testing.T) {
	assert := assert.New(t)
	ctx := context{assert: assert}
	ctx.createBaseFolder()
	ctx.createServer()
	defer ctx.client.Close()
	pathname, token := ctx.startFile("test")
	stream := ctx.startStream(token)
	_, err := os.Stat(path.Join(pathname, "stream"))
	assert.False(os.IsNotExist(err))
	assert.NotEqual(0, len(stream))
}

func TestSynchronizerWriteStream(t *testing.T) {
	assert := assert.New(t)
	ctx := context{assert: assert}
	ctx.createBaseFolder()
	ctx.createServer()
	defer ctx.client.Close()
	pathname, token := ctx.startFile("test")
	stream := ctx.startStream(token)
	ctx.writeStream(stream, []byte("123456"))
	ctx.closeStream(stream)
	data, err := ioutil.ReadFile(path.Join(pathname, "stream"))
	assert.NoError(err)
	assert.Equal([]byte("123456"), data)
}
