package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/service"
)

func MachineServicesStop(_ *connections.Connection, _ *http.Request, params Params) (any, error) {

	svc, err := params.RequireString("service")
	if err != nil {
		return nil, err
	}

	if err := service.PerformAction(svc, service.Stop); err != nil {
		return nil, fmt.Errorf("failed to stop service %q: %w", svc, err)
	}
	return "ok", nil
}
