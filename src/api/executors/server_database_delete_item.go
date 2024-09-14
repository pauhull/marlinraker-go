package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/database"
	"marlinraker/src/marlinraker/connections"
)

type ServerDatabaseDeleteItemResult struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
}

func ServerDatabaseDeleteItem(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	namespace, err := params.RequireString("namespace")
	if err != nil {
		return nil, fmt.Errorf("namespace: %w", err)
	}

	key, err := params.RequireString("key")
	if err != nil {
		return nil, fmt.Errorf("key: %w", err)
	}

	value, err := database.DeleteItem(namespace, key, false)
	if err != nil {
		return nil, fmt.Errorf("failed to delete item: %w", err)
	}

	return ServerDatabaseDeleteItemResult{
		Namespace: namespace,
		Key:       key,
		Value:     value,
	}, nil
}
