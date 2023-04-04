# slow-fuse

Library to simulate slow filesystems by creating a loopback FUSE mount with injected latency


## Usage

### slowfs
This repository comes with a built-in command `slowfs` that can be used to simulate slow FUSE loopback mounts.

```
$ slowfs --help

usage: slowfs [<flags>] [<mountpoint>] [<source>]

A FUSE mount with added latency

Flags:
      --help       Show context-sensitive help (also try --help-long and --help-man).
      --latency=0  Extra latency to add for each listing operation (in milliseconds)
  -v, --verbose    Verbose logs

Args:
  [<mountpoint>]  The location at which source is to be mounted
  [<source>]      Directory to mount
```

### API

#### **Creating a new SlowFUSE mount with 100ms latency**

```go
slowfs, err := slowfuse.New(root, 100*uint64(time.Millisecond))
server, err := slowfs.Mount(mountPoint)
server.Wait()
```

## Caveats
Currently, latency is only injected for the following FUSE operations:
  - getattr
  - statfs
  - open
  - opendir
  - lookup

Create a PR or an issue for more fine-grained latency
