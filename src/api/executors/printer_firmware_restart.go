package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
)

func PrinterFirmwareRestart(*connections.Connection, *http.Request, Params) (any, error) {
	if marlinraker.Printer != nil {
		if err := marlinraker.Printer.Disconnect(); err != nil {
			return nil, fmt.Errorf("failed to disconnect: %w", err)
		}
	}
	go marlinraker.Connect()
	return "ok", nil
}
