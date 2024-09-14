package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/files"
	"marlinraker/src/marlinraker/connections"
)

type ServerFilesZipResult files.ZipAction

func ServerFilesZip(_ *connections.Connection, _ *http.Request, params Params) (any, error) {

	dest, _ := params.RequirePath("dest")

	items, err := params.RequireStringSlice("items")
	if err != nil {
		return nil, err
	}

	storeOnly, _ := params.GetBool("store_only")

	action, err := files.CreateArchive(dest, items, !storeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive: %w", err)
	}
	return action, nil
}
