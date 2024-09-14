package notification

import (
	"encoding/json"
	"fmt"

	"marlinraker/src/marlinraker/connections"
)

type Notification interface {
	Method() string
	Params() []any
}

type rpc struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params,omitempty"`
}

var Testing = false

func Publish(notification Notification) error {
	if Testing {
		return nil
	}

	jsonBytes, err := json.Marshal(rpc{
		JSONRPC: "2.0",
		Method:  notification.Method(),
		Params:  notification.Params(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	for _, connection := range connections.GetConnections() {
		err = connection.WriteText(jsonBytes)
		if err != nil {
			return fmt.Errorf("failed to write to connection: %w", err)
		}
	}
	return nil
}

func Send(connection *connections.Connection, notification Notification) error {
	if Testing {
		return nil
	}

	jsonBytes, err := json.Marshal(rpc{
		JSONRPC: "2.0",
		Method:  notification.Method(),
		Params:  notification.Params(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	err = connection.WriteText(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to write to connection: %w", err)
	}
	return nil
}
