package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/database"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
)

type ServerDatabaseGetItemResult struct {
	Namespace string  `json:"namespace"`
	Key       *string `json:"key"`
	Value     any     `json:"value"`
}

func ServerDatabaseGetItem(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	namespace, err := params.RequireString("namespace")
	if err != nil {
		return nil, fmt.Errorf("namespace: %w", err)
	}

	key, _ := params.GetString("key")

	value, err := database.GetItem(namespace, key, false)
	if err != nil {
		return nil, fmt.Errorf("could not get item %q from namespace %q: %w", key, namespace, err)
	}

	return ServerDatabaseGetItemResult{
		Namespace: namespace,
		Key:       util.StringOrNil(key),
		Value:     value,
	}, nil
}
