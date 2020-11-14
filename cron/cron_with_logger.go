package cron

import (
	"context"
	"time"

	log "github.com/aiyaya188/go-libs/logger"
	"golang.org/x/sync/errgroup"
)

func RunWithLogger(ctx context.Context, logger log.Logger, jobs ...Job) error {
	var g errgroup.Group
	for _, job := range jobs {
		g.Go(runJobFuncWithLogger(ctx, logger, job))
	}
	return g.Wait()
}

func runJobFuncWithLogger(ctx context.Context, logger log.Logger, job Job) func() error {
	return func() error {
		runJobWithLogger(ctx, logger, job)
		return nil
	}
}

func runJobWithLogger(ctx context.Context, logger log.Logger, job Job) {
	logger.Infof("cron %s started", job.Name())
	timer := time.NewTimer(job.Interval())
	defer timer.Stop()
	var timerRead bool
	for {
		select {
		case <-ctx.Done():
			logger.Infof("cron %s exit", job.Name())
			return
		default:
		}

		start := time.Now().UTC()
		logger.Debugf("cron %s job run", job.Name())
		if err := isolate(job.Run)(); err != nil {
			logger.Errorf("cron %s failed: %s", job.Name(), err.Error())
		} else {
			logger.Debugf("cron %s done in %v", job.Name(), time.Since(start))
		}

		if interval := job.Interval() - time.Since(start); interval > 0 {
			if !timer.Stop() && !timerRead {
				<-timer.C
			}
			timer.Reset(interval)
			select {
			case <-ctx.Done():
				logger.Infof("cron %s exit", job.Name())
				return
			case <-timer.C:
				timerRead = true
				continue
			}
		}
	}
}
