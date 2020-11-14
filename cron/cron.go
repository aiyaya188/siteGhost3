package cron

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	log "github.com/aiyaya188/go-libs/logger"
	"golang.org/x/sync/errgroup"
)

type (
	Job interface {
		Name() string
		Interval() time.Duration
		Run() error
	}
)

func Run(jobs ...Job) error {
	var g errgroup.Group
	for _, job := range jobs {
		g.Go(runJobFunc(job))
	}
	return g.Wait()
}

func runJobFunc(job Job) func() error {
	return func() error {
		runJob(job)
		return nil
	}
}

func runJob(job Job) {
	for {
		start := time.Now().UTC()

		//log.Tracef("%s started", job.Name())
		if err := isolate(job.Run)(); err != nil {
			log.Errorf("%s failed: %s", job.Name(), err.Error())
		} else {
			//log.Tracef("%s done in %v", job.Name(), time.Since(start))
		}

		interval := job.Interval() - time.Since(start)
		if interval > 0 {
			time.Sleep(interval)
		}
	}
}

func isolate(f func() error) func() error {
	return func() error {
		defer func() {
			if e := recover(); e != nil {
				err := errors.New(fmt.Sprint(e))
				trace := make([]byte, 4096)
				count := runtime.Stack(trace, true)
				fmt.Fprintf(os.Stderr, "Recover from panic: %s\n", err)
				fmt.Fprintf(os.Stderr, "Stack of %d bytes:\n%s\n", count, trace)
			}
		}()
		return f()
	}
}
