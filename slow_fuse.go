package slowfuse

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
)

type SlowFUSE struct {
	fs.LoopbackNode
	latency  time.Duration
	rootPath string
}

var _ = (fs.NodeStatfser)((*SlowFUSE)(nil))
var _ = (fs.NodeStatfser)((*SlowFUSE)(nil))
var _ = (fs.NodeGetattrer)((*SlowFUSE)(nil))
var _ = (fs.NodeGetxattrer)((*SlowFUSE)(nil))
var _ = (fs.NodeSetxattrer)((*SlowFUSE)(nil))
var _ = (fs.NodeRemovexattrer)((*SlowFUSE)(nil))
var _ = (fs.NodeListxattrer)((*SlowFUSE)(nil))
var _ = (fs.NodeReadlinker)((*SlowFUSE)(nil))
var _ = (fs.NodeOpener)((*SlowFUSE)(nil))
var _ = (fs.NodeCopyFileRanger)((*SlowFUSE)(nil))
var _ = (fs.NodeLookuper)((*SlowFUSE)(nil))
var _ = (fs.NodeOpendirer)((*SlowFUSE)(nil))
var _ = (fs.NodeReaddirer)((*SlowFUSE)(nil))
var _ = (fs.NodeMkdirer)((*SlowFUSE)(nil))
var _ = (fs.NodeMknoder)((*SlowFUSE)(nil))
var _ = (fs.NodeLinker)((*SlowFUSE)(nil))
var _ = (fs.NodeSymlinker)((*SlowFUSE)(nil))
var _ = (fs.NodeUnlinker)((*SlowFUSE)(nil))
var _ = (fs.NodeRmdirer)((*SlowFUSE)(nil))
var _ = (fs.NodeRenamer)((*SlowFUSE)(nil))

func New(rootPath string, latencyNS uint64) (*SlowFUSE, error) {
	var st syscall.Stat_t
	err := syscall.Stat(rootPath, &st)
	if err != nil {
		return nil, err
	}

	root := &fs.LoopbackRoot{
		Path: rootPath,
		Dev:  uint64(st.Dev),
	}

	ret := &SlowFUSE{
		LoopbackNode: fs.LoopbackNode{
			RootData: root,
		},
		latency:  time.Duration(latencyNS),
		rootPath: rootPath,
	}
	log.Debugf("Using latency of %vms", ret.latency.Milliseconds())
	return ret, nil
}

func (s *SlowFUSE) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Debugf("Called Getattr")
	time.Sleep(s.latency)
	ret := s.LoopbackNode.Getattr(ctx, f, out)
	return ret
}

func (s *SlowFUSE) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	log.Debugf("Called Statfs")
	time.Sleep(s.latency)
	return s.LoopbackNode.Statfs(ctx, out)
}

func (s *SlowFUSE) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Debugf("Called Open")
	time.Sleep(s.latency)
	return s.LoopbackNode.Open(ctx, flags)
}

func (s *SlowFUSE) Opendir(ctx context.Context) syscall.Errno {
	log.Debugf("Called Opendir")
	time.Sleep(s.latency)
	return s.LoopbackNode.Opendir(ctx)
}

func (s *SlowFUSE) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	log.Debugf("Called Lookup. name=%v", name)
	time.Sleep(s.latency)
	return s.LoopbackNode.Lookup(ctx, name, out)
}

func (s *SlowFUSE) Mount(mountPoint string) (*fuse.Server, error) {
	timeout := 30 * time.Second

	opts := &fs.Options{
		AttrTimeout:  &timeout,
		EntryTimeout: &timeout,
		MountOptions: fuse.MountOptions{},
	}
	opts.MountOptions.Options = append(opts.MountOptions.Options, "default_permissions")
	opts.MountOptions.Options = append(opts.MountOptions.Options, "ro")
	opts.MountOptions.Options = append(opts.MountOptions.Options, fmt.Sprintf("fsname=%s", s.rootPath))
	opts.MountOptions.Name = "loopback"
	opts.NullPermissions = true

	server, err := fs.Mount(mountPoint, s, opts)
	return server, err
}
