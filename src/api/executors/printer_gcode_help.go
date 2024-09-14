package executors

import (
	"net/http"

	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
)

type PrinterGcodeHelpResult map[string]string

func PrinterGcodeHelp(*connections.Connection, *http.Request, Params) (any, error) {
	if marlinraker.Printer == nil {
		return nil, util.NewError(500, "printer is not online")
	}

	help := make(PrinterGcodeHelpResult)
	for name, macro := range marlinraker.Printer.MacroManager.Macros {
		help[name] = macro.Description()
	}
	return help, nil
}
