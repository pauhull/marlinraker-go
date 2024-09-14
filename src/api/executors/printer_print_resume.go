package executors

import (
	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
	"net/http"
)

func PrinterPrintResume(*connections.Connection, *http.Request, Params) (any, error) {
	if marlinraker.Printer == nil {
		return nil, util.NewError(500, "printer is not online")
	}
	<-marlinraker.Printer.MainExecutorContext().QueueGcode("RESUME", true)
	return "ok", nil
}
