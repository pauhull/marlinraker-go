package executors

import (
	"net/http"

	"marlinraker/src/database"
	"marlinraker/src/marlinraker/connections"
)

type ServerDatabaseListResult struct {
	Namespaces []string `json:"namespaces"`
}

func ServerDatabaseList(*connections.Connection, *http.Request, Params) (any, error) {
	return ServerDatabaseListResult{
		Namespaces: database.ListNamespaces(),
	}, nil
}
