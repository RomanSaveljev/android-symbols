package processor

import (
	"github.com/RomanSaveljev/android-symbols/transmitter/chunker"
	"github.com/RomanSaveljev/android-symbols/transmitter/compressor"
	"github.com/RomanSaveljev/android-symbols/transmitter/encoder"
	"github.com/RomanSaveljev/android-symbols/transmitter/receiver"
	rxapp "github.com/RomanSaveljev/android-symbols/receiver/src/lib"
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func ProcessFileSync(contents []byte, rcv receiver.Receiver, progress chan<- int) (err error) {
	encoder := encoder.NewEncoder(rcv)
	chunker := chunker.NewChunker(encoder, rcv)
	compressor := compressor.NewCompressor(chunker, rcv)
	defer compressor.Close()

	var piece []byte
	for err == nil && len(contents) > 0 {
		available := min(len(contents), rxapp.CHUNK_SIZE)
		piece, contents = contents[0:available], contents[available:]
		_, err = compressor.Write(piece)
		progress <- available
	}

	return
}
