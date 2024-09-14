package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"marlinraker/src/api/executors"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/printer_objects"
	"marlinraker/src/util"
)

type RPC struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type RPCRequest struct {
	RPC
	Method string           `json:"method"`
	Params executors.Params `json:"params"`
}

type RPCResultResponse struct {
	RPC
	Result any `json:"result"`
}

type RPCErrorResponse struct {
	RPC
	Error Error `json:"error"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

func handleSocket(writer http.ResponseWriter, request *http.Request) error {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	connection := connections.RegisterConnection(socket)

	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			break
		}
		log.Debugf("recv: %s", string(message))

		var request RPCRequest
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Errorf("Failed to unmarshal request: %v", err)
			continue
		}

		executor := socketExecutors[request.Method]
		if executor == nil {
			log.Errorf("No executor found for %s", request.Method)
			err = connection.WriteJSON(&RPCErrorResponse{
				Error: Error{404, "Not Found"},
				RPC:   request.RPC,
			})
			if err != nil {
				log.Errorf("Failed to send response: %v", err)
			}
			continue
		}

		result, err := executor(connection, nil, request.Params)
		if err != nil {
			log.Errorf("Error while executing %s: %v", request.Method, err)
			code := 500
			var executorError *util.ExecutorError
			if errors.As(err, &executorError) {
				code = executorError.Code
			}
			err = connection.WriteJSON(&RPCErrorResponse{
				Error: Error{code, err.Error()},
				RPC:   request.RPC,
			})
			if err != nil {
				log.Errorf("Failed to send response: %v", err)
			}
			continue
		}

		err = connection.WriteJSON(&RPCResultResponse{request.RPC, result})
		if err != nil {
			log.Errorf("Failed to send response: %v", err)
		}
	}

	connections.UnregisterConnection(connection)
	printer_objects.Unsubscribe(connection)
	//nolint:nilerr
	return nil
}
