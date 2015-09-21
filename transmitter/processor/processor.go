package processor

import (
	rxapp "github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/encoder"
	"github.com/RomanSaveljev/android-symbols/transmitter/chunker"
	"github.com/RomanSaveljev/android-symbols/transmitter/compressor"
	"github.com/RomanSaveljev/android-symbols/transmitter/receiver"
	"io"
)

func ProcessFileSync(reader io.Reader, rcv receiver.Receiver) (err error) {
	encoder := encoder.NewEncoder(rcv)
	chunker := chunker.NewChunker(encoder, rcv)
	compressor := compressor.NewCompressor(chunker, rcv)
	defer compressor.Close()

	buffer := make([]byte, rxapp.CHUNK_SIZE)
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
