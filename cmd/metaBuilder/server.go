package main

import (
	"context"
	"os"
	"os/signal"
	"siteGhost3/cron"

	log "github.com/aiyaya188/go-libs/logger"
	"golang.org/x/sync/errgroup"
)

type Server struct {
}

func (s *Server) run() error {
	jobs := []cron.Job{
		//daemon.NewFtGqGmaes(),
	}
	go cron.Run(jobs...)
	//daemon.DetectLoginAcount()
	wg, _ := errgroup.WithContext(context.Background())
	wg.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill)
		<-quit
		log.Info("Shutdown Server ...")
		return nil
	})
	wg.Go(func() error {
		DoLoopJob()
		return nil
	})
	return wg.Wait()
}
