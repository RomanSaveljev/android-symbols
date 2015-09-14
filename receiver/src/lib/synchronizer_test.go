package receiver

import (
    "testing"
    "io"
    "io/ioutil"
    "github.com/stretchr/testify/assert"
    "net/rpc"
    "path"
    //"log"
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
	t *testing.T
	in pipePair
	out pipePair
	client *rpc.Client
	baseFolder string
	pathname string
	token string
}

func (this *context) createServer() {
	rpc.DefaultServer = rpc.NewServer()
	this.in.reader, this.out.writer = io.Pipe()
	this.out.reader, this.in.writer = io.Pipe()
	var err error
	this.baseFolder, err = ioutil.TempDir("", "test_")
	assert.Equal(this.t, nil, err)
	go RunSynchronizerService(&this.out)
	this.client = rpc.NewClient(&this.in)
}

func (this *context) startFile(name string) {
	rpc.DefaultServer = rpc.NewServer()
	this.pathname = path.Join(this.baseFolder, name)
	err := this.client.Call("Synchronizer.StartFile", this.pathname, &this.token)
	assert.Equal(this.t, nil, err)
	assert.NotEqual(this.t, 0, len(this.token))
}

func TestSynchronizerStartFile(t *testing.T) {
	var ctx context
	ctx.t = t
	ctx.createServer()
	ctx.startFile("test")
	_, err := os.Stat(ctx.pathname)
	assert.False(t, os.IsNotExist(err))
	ctx.client.Close()
}

func TestSynchronizerFileSignatures(t *testing.T) {
	var ctx context
	ctx.t = t
	ctx.createServer()
	ctx.startFile("test")
}