package receiver

const CHUNK_SIZE = 8*4096

type Chunk struct {
	Signature
	Data [CHUNK_SIZE]byte
}