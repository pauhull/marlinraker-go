package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesPostDirectoryResult files.DirectoryAction

func ServerFilesPostDirectory(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	path, err := params.RequirePath("path")
	if err != nil {
		return nil, fmt.Errorf("path: %w", err)
	}

	action, err := files.CreateDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	return action, nil
}
