package executors

import (
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesRootsResult []files.FileRoot

func ServerFilesRoots(*connections.Connection, *http.Request, Params) (any, error) {
	return files.FileRoots, nil
}
