package executors

import (
	"net/http"

	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
)

func PrinterPrintCancel(*connections.Connection, *http.Request, Params) (any, error) {
	if marlinraker.Printer == nil {
		return nil, util.NewError(500, "printer is not online")
	}
	<-marlinraker.Printer.MainExecutorContext().QueueGcode("CANCEL_PRINT", true)
	return "ok", nil
}
