package receiver

import (
	"github.com/RomanSaveljev/android-symbols/shared/src/shared"
	"net/rpc"
	"os"
	"path"
)

type File struct {
	pathName string
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

func (this *File) Signatures(dummy int, signatures *shared.Signatures) (err error) {
	if file, err := os.Open(this.pathName); err == nil {
		if entries, err := file.Readdirnames(0); err == nil {
			for i := 0; i < len(entries) && err == nil; i++ {
				err = this.goIntoRollingFolder(entries[i], signatures)
			}
		}
	}
	return err
}

func (this *File) StartStream(dummy int, token *string) (err error) {
	*token = path.Join(this.pathName, "stream")
	stream, err := NewStream(*token)
	if err == nil {
		err = rpc.RegisterName(*token, stream)
	}
	return err
}
