package processor

import (
	_ "github.com/RomanSaveljev/android-symbols/receiver/src/lib"
	"github.com/RomanSaveljev/android-symbols/transmitter/encoder"
	"github.com/RomanSaveljev/android-symbols/transmitter/chunker"
	"github.com/RomanSaveljev/android-symbols/transmitter/compressor"
	"github.com/RomanSaveljev/android-symbols/transmitter/receiver"
	"github.com/edsrzf/mmap-go"
	_ "io"
	"os"
)

func ProcessFileSync(file *os.File, rcv receiver.Receiver) (err error) {
	encoder := encoder.NewEncoder(rcv)
	chunker := chunker.NewChunker(encoder, rcv)
	compressor := compressor.NewCompressor(chunker, rcv)
	defer compressor.Close()

	mm, err := mmap.Map(file, mmap.RDONLY, 0)
	defer mm.Unmap()
	if err == nil {
		_, err = compressor.Write(mm)
	}
	/*
	_, err = compressor.Write(buffer[:n])
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
	*/
	return
}
