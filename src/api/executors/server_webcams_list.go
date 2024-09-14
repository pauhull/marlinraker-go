package executors

import (
	"net/http"

	"marlinraker/src/marlinraker/connections"
)

type ServerWebcamsListResult struct {
	// TODO
	Webcams []any `json:"webcams"`
}

func ServerWebcamsList(*connections.Connection, *http.Request, Params) (any, error) {
	return ServerWebcamsListResult{
		Webcams: []any{},
	}, nil
}
