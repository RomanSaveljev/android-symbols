package receiver

import (

)

type Chunk struct {
	Signature
	data [4096]byte
}