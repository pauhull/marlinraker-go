package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesGetDirectoryResult files.DirectoryInfo

func ServerFilesGetDirectory(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	path, err := params.RequirePath("path")
	if err != nil {
		return nil, fmt.Errorf("path: %w", err)
	}

	extended, _ := params.GetBool("extended")

	action, err := files.GetDirInfo(path, extended)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory info: %w", err)
	}
	return action, nil
}
