package executors

import (
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/marlinraker/temp_store"
)

type ServerTemperatureStoreResult temp_store.TempStore

func ServerTemperatureStore(*connections.Connection, *http.Request, Params) (any, error) {
	return temp_store.GetStore(), nil
}
