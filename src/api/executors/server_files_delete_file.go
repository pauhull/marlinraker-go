package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesDeleteFileResult files.FileDeleteAction

func ServerFilesDeleteFile(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	path, err := params.RequirePath("path")
	if err != nil {
		return nil, fmt.Errorf("path: %w", err)
	}

	action, err := files.DeleteFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not delete file %q: %w", path, err)
	}
	return action, nil
}
