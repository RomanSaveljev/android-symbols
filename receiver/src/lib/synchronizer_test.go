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
    "fmt"
    "github.com/RomanSaveljev/android-symbols/shared/src/shared"
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
}

func (this *context) createBaseFolder() {
	var err error
	this.baseFolder, err = ioutil.TempDir("", "test_")
	assert.Equal(this.t, nil, err)	
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
	assert.Equal(this.t, nil, err)
	return
}

func (this *context) getSignatures(token string) *shared.Signatures {
	signatures := shared.NewSignatures()
	err := this.client.Call(fmt.Sprint(token, ".Signatures"), signatures, signatures)
	assert.Nil(this.t, err)
	return signatures
}

func TestSynchronizerStartFile(t *testing.T) {
	var ctx context
	ctx.t = t
	ctx.createBaseFolder()
	ctx.createServer()
	pathname, token := ctx.startFile("test")
	_, err := os.Stat(pathname)
	assert.False(t, os.IsNotExist(err))
	assert.NotEqual(t, 0, len(token))
	ctx.client.Close()
}

func TestSynchronizerFileSignatures(t *testing.T) {
	var ctx context
	ctx.t = t
	ctx.createBaseFolder()
	sigPath := path.Join(ctx.baseFolder, "test", "ffffeeee")
	err := os.MkdirAll(sigPath, os.ModeDir|os.ModePerm)
	assert.Equal(t, nil, err)
	err = ioutil.WriteFile(path.Join(sigPath, "6a900ff954cc07f2ea7b0fa2b862c480"), []byte("123"), os.ModePerm)
	assert.Equal(t, nil, err)
	err = ioutil.WriteFile(path.Join(sigPath, "449448c6292cfab8fd2038c65a2727a9"), []byte("456"), os.ModePerm)
	assert.Nil(t, err)
	ctx.createServer()
	_, token := ctx.startFile("test")
	signatures := ctx.getSignatures(token)
	strong := signatures.Get("ffffeeee")
	assert.Equal(t, 2, len(strong))
	assert.Contains(t, strong, "6a900ff954cc07f2ea7b0fa2b862c480")
	assert.Contains(t, strong, "449448c6292cfab8fd2038c65a2727a9")
	ctx.client.Close()
}