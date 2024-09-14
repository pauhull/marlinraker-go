package executors

import (
	"net/http"

	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/util"
)

type ServerConnectionIdentifyResult struct {
	ConnectionID int `json:"connection_id"`
}

func ServerConnectionIdentify(connection *connections.Connection, _ *http.Request, params Params) (any, error) {

	if connection.Identified {
		return nil, util.NewError(400, "connection already identified")
	}

	clientName, err := params.RequireString("client_name")
	if err != nil {
		return nil, err
	}

	version, err := params.RequireString("version")
	if err != nil {
		return nil, err
	}

	clientType, err := params.RequireString("type")
	if err != nil {
		return nil, err
	}

	url, err := params.RequireString("url")
	if err != nil {
		return nil, err
	}

	connection.ClientName, connection.Version, connection.ClientType, connection.URL =
		clientName, version, clientType, url

	return ServerConnectionIdentifyResult{connection.ID}, nil
}
