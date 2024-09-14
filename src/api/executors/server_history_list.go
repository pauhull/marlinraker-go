// TODO

package executors

import (
	"net/http"

	"marlinraker/src/marlinraker/connections"
)

type ServerHistoryListResult struct {
	Count int   `json:"count"`
	Jobs  []any `json:"jobs"`
}

func ServerHistoryList(*connections.Connection, *http.Request, Params) (any, error) {
	return ServerHistoryListResult{
		Count: 0, Jobs: make([]any, 0),
	}, nil
}
