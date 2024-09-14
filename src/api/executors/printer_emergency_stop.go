package executors

import (
	"net/http"

	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
)

type PrinterEmergencyStopResult string

func PrinterEmergencyStop(*connections.Connection, *http.Request, Params) (any, error) {
	marlinraker.Printer.EmergencyStop()
	return "ok", nil
}
