package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesMoveResult files.MoveAction

func ServerFilesMove(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	source, err := params.RequirePath("source")
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}

	dest, err := params.RequirePath("dest")
	if err != nil {
		return nil, fmt.Errorf("dest: %w", err)
	}

	action, err := files.Move(source, dest)
	if err != nil {
		return nil, fmt.Errorf("failed to move file: %w", err)
	}
	return action, nil
}
