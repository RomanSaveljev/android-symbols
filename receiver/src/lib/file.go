package receiver

import (
	"github.com/RomanSaveljev/android-symbols/shared/src/shared"
	"io"
	"net/rpc"
	"os"
	"path"
)

type File struct {
	pathName    string
	rollingDirs *[]string
	strongFiles *[]string
}

func NewFile(pathName string) (*File, error) {
	var f File
	f.pathName = pathName
	err := os.MkdirAll(pathName, os.ModeDir|os.ModePerm)
	return &f, err
}

func (this *File) goIntoRollingFolder(entry string, signatures *shared.Signatures) (err error) {
	if file, err := os.Open(path.Join(this.pathName, entry)); err == nil {
		if entries, err := file.Readdirnames(0); err == nil {
			for i := 0; i < len(entries); i++ {
				signatures.Add(entry, entries[i])
			}
		}
	}
	return err
}

func (this *File) NextSignature(dummy int, signature *Signature) (err error) {
	if this.rollingDirs == nil {
		// lets begin
		if file, err := os.Open(this.pathName); err == nil {
			entries, err := file.Readdirnames(0)
			_ = file.Close()
			if err == nil {
				this.rollingDirs = &entries
				return this.NextSignature(dummy, signature)
			}
		}
	} else if len(*this.rollingDirs) == 0 {
		// all rolling directories scanned
		err = io.EOF
	} else if this.strongFiles == nil {
		// entered a rolling directory - scan strong signatures
		entry := (*this.rollingDirs)[0]
		if file, err := os.Open(path.Join(this.pathName, entry)); err == nil {
			entries, err := file.Readdirnames(0)
			file.Close()
			if err == nil {
				this.strongFiles = &entries
				next := (*this.rollingDirs)[1:]
				this.rollingDirs = &next
				return this.NextSignature(dummy, signature)
			}
		}
	} else if len(*this.strongFiles) == 0 {
		// all strong signatures listed
		this.strongFiles = nil
		return this.NextSignature(dummy, signature)
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
