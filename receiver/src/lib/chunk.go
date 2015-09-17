package receiver

const CHUNK_SIZE = 4096

type Chunk struct {
	Signature
	data [CHUNK_SIZE]byte
}