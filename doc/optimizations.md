# Optimizations to improve the performance

Always use profiler!

## New file mode

When we know the file did not exist in the remote side, we can pump it directly without checking the
signatures, because we know there is none. The **compressor** always sends byte after byte into the **chunker**. When its buffer gets full, the **chunker** has to cound rolling and strong signatures anyway. Current list of signatures will be updating. If the chunk with the same content has been transmitted already, then it only records a signature to the stream. Some traffic will be saved.  