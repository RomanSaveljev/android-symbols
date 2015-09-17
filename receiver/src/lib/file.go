package receiver

//go:generate $GOPATH/bin/mockgen -source file.go -package receiver -destination mock_file_system_worker_test.go

import (
	"io"
	"log"
	"net/rpc"
	"os"
	"path"
)

type fileSystemWorker interface {
	Readdirnames(at string) ([]string, error)
	MkdirAll(pathName string) error
}

type realFileSystemWorker struct {
	fileSystemWorker
}

func (this realFileSystemWorker) Readdirnames(at string) (entries []string, err error) {
	var file *os.File
	if file, err = os.Open(at); err == nil {
		defer file.Close()
		entries, err = file.Readdirnames(0)
	}
	return
}

func (this realFileSystemWorker) MkdirAll(pathName string) error {
	return os.MkdirAll(pathName, os.ModeDir|os.ModePerm)
}

type File struct {
	pathName    string
	rollingDirs *[]string
	strongFiles *[]string
	worker      fileSystemWorker
}

func NewFile(pathName string) (*File, error) {
	return newFileInjected(pathName, realFileSystemWorker{})
}

func newFileInjected(pathName string, worker fileSystemWorker) (*File, error) {
	f := File{pathName: pathName, worker: worker}
	err := worker.MkdirAll(pathName)
	return &f, err
}

func (this *File) NextSignature(dummy int, signature *Signature) (err error) {
	*signature, err = this.nextSignature()
	return
}

func (this *File) nextSignature() (signature Signature, err error) {
	if this.rollingDirs == nil {
		// lets begin
		log.Println("beginning")
		if entries, err := this.worker.Readdirnames(this.pathName); err == nil {
			log.Printf("%v", entries)
			this.rollingDirs = &entries
			return this.nextSignature()
		}
	} else if len(*this.rollingDirs) == 0 {
		// all rolling directories scanned
		err = io.EOF
	} else if this.strongFiles == nil {
		// entered a rolling directory - scan strong signatures
		log.Println("Fetching strong")
		entry := (*this.rollingDirs)[0]
		if entries, err := this.worker.Readdirnames(path.Join(this.pathName, entry)); err == nil {
			this.strongFiles = &entries
			return this.nextSignature()
		}
	} else if len(*this.strongFiles) == 0 {
		// all strong signatures listed
		this.strongFiles = nil
		next := (*this.rollingDirs)[1:]
		this.rollingDirs = &next
		return this.nextSignature()
	} else {
		signature.Rolling = (*this.rollingDirs)[0]
		signature.Strong = (*this.strongFiles)[0]
		next := (*this.strongFiles)[1:]
		this.strongFiles = &next
	}
	return
}

func (this *File) StartStream(dummy int, token *string) (err error) {
	*token = path.Join(this.pathName, "stream")
	stream, err := NewStream(*token)
	if err == nil {
		err = rpc.RegisterName(*token, stream)
	}
	return err
}
