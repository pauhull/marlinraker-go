package executors

import (
	"net/http"
	"strconv"

	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
)

func PrinterPrintStart(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	if marlinraker.Printer == nil {
		return nil, util.NewError(500, "printer is not online")
	}

	fileName, err := params.RequireString("filename")
	if err != nil {
		return nil, err
	}

	<-marlinraker.Printer.MainExecutorContext().QueueGcode("SDCARD_PRINT_FILE FILENAME="+strconv.Quote(fileName), true)
	return "ok", nil
}
