package main

import (
	"dishcovery/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	// Session  *scs.SessionManager
	DB        *gorm.DB
	InfoLog   *zap.SugaredLogger
	ErrorLog  *zap.SugaredLogger
	Models    data.Models
	Wait      *sync.WaitGroup
	ErrorChan chan error
	DoneChan  chan bool
}

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Infof("Starting web server at %s...", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// perform any cleanup tasks
	app.InfoLog.Info("would run cleanup tasks...")

	// block until waitgroup is empty
	app.Wait.Wait()
	app.DoneChan <- true

	app.InfoLog.Info("closing channels and shutting down application...")
	close(app.ErrorChan)
	close(app.DoneChan)
}

func (app *Config) listenForErrors() {
	for {
		select {
		case err := <-app.ErrorChan:
			//error handling
			_ = err
		case <-app.DoneChan:
			return
		}
	}
}
