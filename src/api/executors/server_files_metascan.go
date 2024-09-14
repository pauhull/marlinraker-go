package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesMetascanResult *files.Metadata

func ServerFilesMetascan(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	fileName, err := params.RequirePath("filename")
	if err != nil {
		return nil, fmt.Errorf("filename: %w", err)
	}

	_ = files.RemoveMetadata(fileName)

	action, err := files.LoadOrScanMetadata(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to scan metadata: %w", err)
	}
	return action, nil
}
