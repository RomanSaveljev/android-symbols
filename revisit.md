## New incarnation of "android-symbols"

The goal of the project is to offer a version content management system optimized
for binary blobs. Merging and rebasing are not primary features. Rather it should
implement a simple snapshot system with branching.

Revisited features:

* Decouple data store from the processing - they have to be in two different
containers
* Extract file(s) from a layered structure on request
* Save file(s) to a layered structure on request


## Technical details

This "simple" construct will send individual files tarred and accept and untar
them at the receiving end:

```
(for FILE in a b c d; do echo 123 && tar c $FILE; done) | (while read LINE; do tar xv -C out;done)
```

This command will mount `tar` archive:

```
archivemount tar.tar /tmp/tar-test/mnt/
```

Individual files then can be split and transferred to the receiving end.

Go has [tar](https://golang.org/pkg/archive/tar) package, which only needs a very
simple [io.Reader](https://golang.org/pkg/io/#Reader) to extract entries from the
archive.

A file chunk is extracted to the memory and then it is compared to the one stored.
If they are different, the new one is saved. Otherwise, nothing happens to avoid
changes in the layer.

## Storage container organization

It is created from [scratch](https://hub.docker.com/_/scratch/). Controlled files
can sit at its root. To update some files we need to:

* Create a container
* Copy modified file to the container
* Commit the container as a new image

Newer version of docker allows to copy from and into container's filesystem.

* User provides a set of files to store
* Each file is processed in its own go-routine
* Each file is broken into chunks of a fixed size
* Each chunk is processed in its own go-routine
* Each chunk is compared with its counterpart inside a container
* Different chunks are copied into the container's filesystem

## Volume performance

20 files in the volume:

```shell
$ time docker create volume-test
bce8baa92469be90f392878ae5e6ea29f00440404317ddb6fff9f8ff72053a6c

real	0m0.357s
user	0m0.020s
sys	0m0.020s
$ time docker create --volume /opt/volume volume-test
c95c03d9102c5f4304d6fb40a8b1fd718ea3fa7b63e59ede67c877eca0703788

real	0m0.442s
user	0m0.018s
sys	0m0.018s
```

20000 files in the volume:

```shell
$ time docker create volume-test
ff6df8e30630e0f34852d9bf9a0103073c76bcc6bb32daf2dea93cd5ed2d8a4f

real	0m0.437s
user	0m0.017s
sys	0m0.011s
$ time docker create --volume /opt/volume volume-test
d41766844e32a58c28118d681303071e63896a9ffc761ee093095bb7dcb5c8c5

real	0m1.802s
user	0m0.035s
sys	0m0.000s
```

20000 8 kB files:

```shell
$ time docker create volume-test
4eeac0140cdca31f18966aaffe73443a4b5a983d5131243ac1aa20b7e00f282d

real	0m0.624s
user	0m0.013s
sys	0m0.017s
$ time docker create --volume /opt/volume volume-test
37efbf5be63638d7c56ef3f2d73def4f4e6627008c019c9739f18f8c41360cc3

real	0m7.223s
user	0m0.020s
sys	0m0.016s
```

Creating a container with a volume creates a snapshot of the content in that
volume and actually uses disk space. The price to pay seems unreasonable..

To operate on the storage container side we would have to plant some executable
code there. This applies some restrictions on the storage container structure.

Go executable can run without external dependencies (doing syscalls). There should
be a folder with a reserved name to hold the executable. TX and RX will communicate
over STDOUT and STDIN. Multiple parallels message exchanges need to happen at the
same time over the same channel. Need to see what kind of transport can be used...

The exchange protocol should be transport-agnostic. Exchange over network will be
useful for "server push". The [rpc](http://golang.org/pkg/net/rpc) implementation
of Go should be sufficient. On the other hand, simple STDIN-STDOUT duplex should
be enough, because docker client can work well with remote daemons. It will
implement security and authentication as well.

Transmitter and receiver can randomly access their versions of the file. That is
transmitter can support the receiver and send parts of the file upon request.

## Storing files in a layer

A file is a granularity unit in a layer. If one byte modification is applied to
a 1 gigabyte file, then new layer size will be 1 gigabyte.

### Break files into equally sized chunks

Pros:

* One byte change only grows a layer by the size of a chunk
* Effort to retrieve a file does not depend on the number of versions

Cons:

* Insert one byte at the start of the file - all chunks have to be rebuilt
* Have to reassemble a file every time, even if it was never changed

### Store diffs incrementally

With every change a layer will store a patch

Pros:

* Layer overhead closely reflects the amount of modified data
* Files without modifications are recovered quicker
* Merging of diffs and other optimizations are possible

Cons:

* Files with many versions are retrieved slower - can be managed by squashing
diffs
* No good libraries for binary diff and patch
* Need to apply all patches to decide that file has changed


### Keep the rsync products and do not assemble the file

If the transmission between transmitter and receiver implements rsync protocol,
then we would not need to save an assembled file on the layer. Rsync transmits
roughly a difference, which should be saved to a layer.

When the file is transmitted for the first time it is simple broke into chunks on
the receiving end. For every chunk signatures are calculated and used to name to
chunk files. Chunks are arranged into a tree structure:

```
root
--hash
----quick-signature
------slow-signature
--stream
```

In the rsync algorithm received chunks are used to build a common dictionary. The
perfect dictionary is a complete set of all possible chunks. The `stream` in the
above figure is a sequence of literal bytes interspersed with chunk identities.

The stream file is always updated between layers, but chunks collection is only
extended. The more literal bytes we have in a stream file the bigger it is and
the more expensive a layer is.

If there is a continuous sequence of literal bytes as big as an agreed chunk size,
it can be optimized out. The sequence may be made into chunk and replaced by its
identification in the stream file. This uniform procedure will both work for
updating an existing file and transmitting its first version.

Stream file is ASCII85 encoded for simple handling. Sequences of literal bytes
are separated by new line (`\n`) character. A line may either represent a literal
bytes sequence or reference a dictionary chunk. To distinguish the latter, the line
should start with horizontal tabulation (`\t`) character. Chunk data is stored in
binary.

## User interface

User interface is reprensented by CLI, which roughly reflects GIT CLI design. User
updates some files in his/her working area. After that the modified files should
be staged and then the change is committed.

There is no need to try hard to make it work like GIT.

### Working area updates

TBD

### Staging

Modified files are staged by saving them to the storage container. Storage
container has local changes. Committing a storage container preserves it and
then it can be destroyed.
