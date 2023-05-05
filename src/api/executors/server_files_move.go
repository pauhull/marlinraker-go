package executors

import (
	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
	"net/http"
)

type ServerFilesMoveResult files.MoveAction

func ServerFilesMove(_ *connections.Connection, _ *http.Request, params Params) (any, error) {

	source, exists := params.GetString("source")
	if !exists {
		return nil, util.NewError("source param is required", 400)
	}

	dest, exists := params.GetString("dest")
	if !exists {
		return nil, util.NewError("dest param is required", 400)
	}

	return files.Move(source, dest)
}
