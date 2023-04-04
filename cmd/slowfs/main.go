package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	slowfuse "github.com/gurupras/slow-fuse"
)

var (
	app = kingpin.New("slowfs", "A FUSE mount with added latency")

	latency    = app.Flag("latency", "Extra latency to add for each listing operation (in milliseconds)").Default("0").Uint64()
	mountpoint = app.Arg("mountpoint", "The location at which source is to be mounted").String()
	source     = app.Arg("source", "Directory to mount").String()
	verbose    = app.Flag("verbose", "Verbose logs").Default("false").Short('v').Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	slowfs, err := slowfuse.New(*source, *latency*uint64(time.Millisecond))
	if err != nil {
		log.Fatalf("Failed to create new SlowFUSE: %v", err)
	}
	server, err := slowfs.Mount(*mountpoint)
	if err != nil {
		log.Fatalf("Failed to mount: %v", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigc
		server.Unmount()
		os.Exit(-1)
	}()

	server.Wait()
}
