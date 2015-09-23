package processor

import (
	"github.com/RomanSaveljev/android-symbols/transmitter/chunker"
	"github.com/RomanSaveljev/android-symbols/transmitter/compressor"
	"github.com/RomanSaveljev/android-symbols/transmitter/encoder"
	"github.com/RomanSaveljev/android-symbols/transmitter/receiver"
)

func ProcessFileSync(contents []byte, rcv receiver.Receiver) (err error) {
	encoder := encoder.NewEncoder(rcv)
	chunker := chunker.NewChunker(encoder, rcv)
	compressor := compressor.NewCompressor(chunker, rcv)
	defer compressor.Close()

	_, err = compressor.Write(contents)

	return
}
