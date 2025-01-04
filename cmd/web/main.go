package main

import (
	"dishcovery/data"
	"dishcovery/handler/dbHandler"
	"dishcovery/handler/logHandler"
	"sync"

	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

var webPort = "8080"

func main() {

	db := dbHandler.InitDB()

	app := Config{
		DB:        db,
		InfoLog:   logHandler.InitLogger(),
		ErrorLog:  logHandler.InitLogger(),
		Models:    data.New(db),
		Wait:      &sync.WaitGroup{},
		ErrorChan: make(chan error),
		DoneChan:  make(chan bool),
	}
	go app.listenForShutdown()
	go app.listenForErrors()

	app.serve()
}
