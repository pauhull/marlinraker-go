package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesDeleteDirectoryResult files.DirectoryAction

func ServerFilesDeleteDirectory(_ *connections.Connection, _ *http.Request, params Params) (any, error) {

	path, err := params.RequirePath("path")
	if err != nil {
		return nil, fmt.Errorf("path: %w", err)
	}

	force, _ := params.GetBool("force")

	action, err := files.DeleteDir(path, force)
	if err != nil {
		return nil, fmt.Errorf("failed to delete directory: %w", err)
	}
	return action, nil
}
