package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/database"
	"marlinraker/src/marlinraker/connections"
)

type ServerDatabasePostItemResult struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
}

func ServerDatabasePostItem(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	namespace, err := params.RequireString("namespace")
	if err != nil {
		return nil, fmt.Errorf("namespace: %w", err)
	}

	key, err := params.RequireString("key")
	if err != nil {
		return nil, fmt.Errorf("key: %w", err)
	}

	value, err := params.RequireAny("value")
	if err != nil {
		return nil, fmt.Errorf("value: %w", err)
	}

	if _, err = database.PostItem(namespace, key, value, false); err != nil {
		return nil, fmt.Errorf("failed to post item: %w", err)
	}

	return ServerDatabasePostItemResult{
		Namespace: namespace,
		Key:       key,
		Value:     value,
	}, nil
}
