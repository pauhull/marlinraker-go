package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/service"
)

func MachineServicesRestart(_ *connections.Connection, _ *http.Request, params Params) (any, error) {

	svc, err := params.RequireString("service")
	if err != nil {
		return nil, err
	}

	if err := service.PerformAction(svc, service.Restart); err != nil {
		return nil, fmt.Errorf("failed to restart service %q: %w", svc, err)
	}
	return "ok", nil
}
