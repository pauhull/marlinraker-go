package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesMetadataResult *files.Metadata

func ServerFilesMetadata(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	fileName, err := params.RequirePath("filename")
	if err != nil {
		return nil, fmt.Errorf("filename: %w", err)
	}

	meta, err := files.LoadOrScanMetadata(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}
	return meta, nil
}
