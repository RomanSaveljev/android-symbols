# Optimizations to improve the performance

Always use profiler!

## New file mode

When we know the file did not exist in the remote side, we can pump it directly without checking the
signatures, because we know there is none. The **compressor** always sends byte after byte into the **chunker**. When its buffer gets full, the **chunker** has to cound rolling and strong signatures anyway. Current list of signatures will be updating. If the chunk with the same content has been transmitted already, then it only records a signature to the stream. Some traffic will be saved.  

## CountRolling is slower than CountStrong

The main problem is the looping through the slice. The loop alone takes 9495 ns/op. Another issue is conversion to `uint32`. A simple `i++` loop takes as much.

I have tried many things, but so far could not speed it up significantly