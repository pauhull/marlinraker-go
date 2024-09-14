package executors

import (
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/system_info"
)

type MachineProcStatsResult = system_info.ProcStats

func MachineProcStats(*connections.Connection, *http.Request, Params) (any, error) {
	return system_info.GetStats(), nil
}
