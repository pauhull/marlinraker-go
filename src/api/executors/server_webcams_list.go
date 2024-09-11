package executors

import (
	"marlinraker/src/marlinraker/connections"
	"net/http"
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
