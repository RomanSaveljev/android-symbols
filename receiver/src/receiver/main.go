package main 

import (
	"net/rpc"
	"github.com/RomanSaveljev/android-symbols/shared/src/shared"
)

type Service int

func (this *Service) FindSignatures(filePath string, sigs *shared.Signatures) error {
	return nil
}

func main() {
	var tr shared.Transport
	var service Service
	rpc.Register(service)
	rpc.ServeConn(&tr)
}

