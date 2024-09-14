package logger

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"marlinraker/src/files"
)

var (
	logFile afero.File
	mu      = &sync.Mutex{}
)

func SetupLogger() error {

	mu.Lock()
	defer mu.Unlock()

	if err := openLogFile(); err != nil {
		return err
	}

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		NoColors:        true,
	})

	return nil
}

func HandleRotate() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)

	for {
		<-ch
		log.Println("Log rotate requested")
		mu.Lock()
		if err := openLogFile(); err != nil {
			panic(err)
		}
		mu.Unlock()
	}
}

func CloseLogger() error {
	mu.Lock()
	defer mu.Unlock()
	if err := logFile.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %w", err)
	}
	return nil
}

func openLogFile() error {

	var err error
	if logFile != nil {
		if err = logFile.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}

	logFilePath := filepath.Join(files.DataDir, "logs/marlinraker.log")
	logFile, err = files.Fs.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return fmt.Errorf("failed to open %q: %w", logFile, err)
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	return nil
}
