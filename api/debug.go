package api

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func init() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)
	ftrace := filepath.Join(os.TempDir(), "swan-ng-stack-trace.log")

	go func() {
		for range ch {
			f, err := os.OpenFile(ftrace, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Error("write stack trace log file error", err)
				continue
			}
			fmt.Fprint(f, "GOROUTINE\n\n")
			pprof.Lookup("goroutine").WriteTo(f, 2)
			fmt.Fprint(f, "\n\nHEAP\n\n")
			pprof.Lookup("heap").WriteTo(f, 1)
			fmt.Fprint(f, "\n\nTHREADCREATE\n\n")
			pprof.Lookup("threadcreate").WriteTo(f, 1)
			fmt.Fprint(f, "\n\nBLOCK\n\n")
			pprof.Lookup("block").WriteTo(f, 1)
			f.Close()
		}
	}()
}
