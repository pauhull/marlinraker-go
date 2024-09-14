package executors

import (
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
)

type ServerFilesResult = []files.File

func ServerFilesList(_ *connections.Connection, _ *http.Request, params Params) (any, error) {
	root, exists := params.GetString("root")
	if !exists {
		root = "gcodes"
	}
	filesList, err := files.ListFiles(root)
	if err != nil {
		return nil, util.NewErrorf(400, "%v", err)
	}
	return filesList, nil
}
