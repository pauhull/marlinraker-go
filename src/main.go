package main

import (
	_ "embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"marlinraker/src/api"
	"marlinraker/src/database"
	"marlinraker/src/files"
	"marlinraker/src/logger"
	"marlinraker/src/marlinraker"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {

	if os.Geteuid() == 0 {
		fmt.Println("Please do not run this program as root")
		return
	}

	dataDir, err := filepath.Abs("./marlinraker_files")
	if err != nil {
		panic(err)
	}

	files.DataDir = dataDir
	if err := files.Init(); err != nil {
		panic(err)
	}

	logFile, err := logger.SetupLogger(dataDir)
	if err != nil {
		panic(err)
	}
	defer logger.CloseLogger(logFile)

	if err := database.Init(); err != nil {
		panic(err)
	}

	marlinraker.Init(dataDir)

	go api.StartServer()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
	log.Println("Received interrupt, shutting down")
}
