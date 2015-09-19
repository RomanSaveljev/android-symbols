package transmitter

import (
	"github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"io"
)

func ProcessFileSync(reader io.Reader, rcv Receiver) (err error) {
	encoder := NewEncoder(rcv)
	chunker := NewChunker(encoder, rcv)
	compressor := NewCompressor(chunker, rcv)
	defer compressor.Close()

	buffer := make([]byte, receiver.CHUNK_SIZE)
	for err == nil {
		var n int
		if n, err = reader.Read(buffer); err == nil {
			_, err = compressor.Write(buffer[:n])
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
