package executors

import (
	"fmt"
	"net/http"

	"marlinraker/src/auth"
	"marlinraker/src/marlinraker/connections"
)

type AccessOneshotTokenResult string

func AccessOneshotToken(*connections.Connection, *http.Request, Params) (any, error) {
	response, err := auth.GenerateOneshotToken()
	if err != nil {
		return nil, fmt.Errorf("could not generate oneshot token: %w", err)
	}
	return AccessOneshotTokenResult(response), nil
}
