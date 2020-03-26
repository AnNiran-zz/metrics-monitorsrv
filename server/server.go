package server

import (
	"os"
	"context"
	"time"
	"os/signal"
	"github.com/argcv/stork/log"
)

func Run() {
	srv := newHttpSrv()
	if err := srv.start(); err != nil {
		log.Fatalf("Start failed: %v", err)
		return 
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infof("[Contol-C] Get signal: shutdown server ...")
	signal.Reset(os.Interrupt)

	// starting shutting down progress
	log.Infof("Server shutting down")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second)
	
	defer cancel()

	if err := srv.shutdown(ctx); err != nil {
		log.Errorf("Server Shutdown failed: %v", err)
	}
	log.Infof("Server exiting")
}
