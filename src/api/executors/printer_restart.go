package executors

import (
	log "github.com/sirupsen/logrus"
	"marlinraker/src/marlinraker/connections"
	"net/http"
	"syscall"
)

func PrinterRestart(*connections.Connection, *http.Request, Params) (any, error) {
	go func() {
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
			log.Errorf("Failed to restart printer: %v", err)
		}
	}()
	return "ok", nil
}
