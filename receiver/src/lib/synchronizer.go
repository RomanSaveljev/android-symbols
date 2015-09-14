package receiver

import (
	"net/rpc"
	"io"
)

type Synchronizer int

func (this *Synchronizer) StartFile(filePath string, token *string) error {
	*token = filePath
	file, err := NewFile(filePath)
	if err == nil {
		err = rpc.RegisterName(filePath, file)
	}
	return err
}

func RunSynchronizerService(link io.ReadWriteCloser) error {
	err := rpc.Register(new(Synchronizer))
	if err == nil {
		rpc.ServeConn(link)
	}
	return err
}
