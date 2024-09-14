package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/service"
)

func MachineServicesStart(_ *connections.Connection, _ *http.Request, params Params) (any, error) {

	svc, err := params.RequireString("service")
	if err != nil {
		return nil, err
	}

	if err := service.PerformAction(svc, service.Start); err != nil {
		return nil, fmt.Errorf("could not start service %q: %w", svc, err)
	}
	return "ok", nil
}
