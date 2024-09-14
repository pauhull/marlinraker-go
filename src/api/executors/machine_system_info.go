package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/system_info"
)

type MachineSystemInfoResult struct {
	SystemInfo *system_info.SystemInfo `json:"system_info"`
}

func MachineSystemInfo(*connections.Connection, *http.Request, Params) (any, error) {

	info, err := system_info.GetSystemInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}

	return MachineSystemInfoResult{info}, nil
}
