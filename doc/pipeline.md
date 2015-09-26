# Transmitter pipeline description

A pipeline handles synchronization of a single file.

## Retrieving chunks

It all starts with querying existing chunks from the receiver's end. A collection of signatures is built on the transmitter's end.

## Compressing

Compressor receives a file to synchronize and begins to walk it byte after byte. At each offset it calculates the signature of a fixed size chunk. If the exact signature is available from the collection, then only the signature itself is sent and file reading position jumps to the byte after the chunk. If the chunk signature is not known, then the byte at the offset is passed to the chunker.

## Chunker

Literal bytes passed from the compressor are buffered here. When compressor writes a signature instead, the buffer is flushed to the receiver's end. If the buffer grows enough, then it is flushed in the form of a signature. Buffer's data is transmitted as a new chunk and a signature is recorded on the transmitter's side. Compressor will use this new signature to check the following bytes.

## Encoder

Encoder converts literal bytes to [ASCII85](https://en.wikipedia.org/wiki/Ascii85) format and mixes them with signatures to prepare the stream file. Encoder will receive bigger portions of data, because of the buffering in the chunker.

# Transmitter GUI

If the transmitter is connected to TTY, it should render the overall progress information as well as current transmission speed. Every compressor will regularly update the total number of bytes they pushed to a dedicated channel. A channel will be allocated per each file. Data from all channels is collected and published to two channels: progress rendering and speed sampling.
 
# Receiver pipeline

Receiver simply collects the data from the encoder and records it to the stream file.