package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"marlinraker/src/api/executors"
	"marlinraker/src/files"
)

func testHTTP[Result any](t *testing.T, method string, endpoint string, params executors.Params,
	f func(*testing.T, *httptest.ResponseRecorder, *Result, *Error)) {
	t.Run(method+endpoint, func(t *testing.T) {

		urlQuery, err := url.Parse(endpoint)
		if err != nil {
			t.Fatal(err)
		}

		var body io.Reader
		if method == "GET" {
			values := urlQuery.Query()
			for param, value := range params {
				values.Set(param, fmt.Sprintf("%v", value))
			}
			urlQuery.RawQuery = values.Encode()
		} else {
			bodyBytes, err := json.Marshal(params)
			require.NoError(t, err)
			body = bytes.NewReader(bodyBytes)
		}

		request, err := http.NewRequestWithContext(context.Background(), method, urlQuery.String(), body)
		require.NoError(t, err)
		if body != nil {
			request.Header.Add("Content-Type", "application/json")
		}

		recorder := httptest.NewRecorder()
		HTTPHandler{}.ServeHTTP(recorder, request)

		var errorResponse ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Fatal(err)
			return
		}
		if errorResponse.Error.Code > 0 {
			f(t, recorder, nil, &errorResponse.Error)
			return
		}

		var response struct {
			Result Result `json:"result"`
		}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
			return
		}
		f(t, recorder, &response.Result, nil)
	})
}

func testSocket[Result any](t *testing.T, method string, params executors.Params, f func(*testing.T, *Result, *Error)) {
	t.Run(method, func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if err := handleSocket(writer, request); err != nil {
				t.Fatal(err)
			}
		}))

		defer server.Close()
		socketURL := fmt.Sprintf("ws%s", server.URL[4:])

		//nolint:bodyclose
		socket, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = socket.Close()
		}()

		request := RPCRequest{
			RPC: RPC{
				JSONRPC: "2.0",
				ID:      0,
			},
			Method: method,
			Params: params,
		}

		{
			data, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}

			err = socket.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				t.Fatal(err)
			}
		}

		{
			_, data, err := socket.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}

			var errorResponse RPCErrorResponse
			err = json.Unmarshal(data, &errorResponse)
			if err != nil {
				t.Fatal(err)
			}
			if errorResponse.Error.Code > 0 {
				f(t, nil, &errorResponse.Error)
				return
			}

			var response struct {
				Result Result `json:"result"`
			}
			err = json.Unmarshal(data, &response)
			if err != nil {
				t.Fatal(err)
			}
			f(t, &response.Result, nil)
		}
	})
}

func testAll[Result any](t *testing.T, rpcMethod string, httpMethod string, httpURL string, params executors.Params,
	f func(*testing.T, *httptest.ResponseRecorder, *Result, *Error)) {

	testSocket(t, rpcMethod, params, func(t *testing.T, result *Result, error *Error) {
		f(t, nil, result, error)
	})
	testHTTP(t, httpMethod, httpURL, params, f)
}

func makeConnection(t *testing.T) (*websocket.Conn, int) {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if err := handleSocket(writer, request); err != nil {
			t.Fatal(err)
		}
	}))

	defer server.Close()
	socketURL := fmt.Sprintf("ws%s", server.URL[4:])

	//nolint:bodyclose
	socket, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	require.NoError(t, err)

	err = socket.WriteJSON(RPCRequest{
		RPC:    RPC{JSONRPC: "2.0", ID: 42},
		Method: "server.connection.identify",
		Params: executors.Params{"client_name": "", "version": "", "type": "", "url": ""},
	})
	if err != nil {
		t.Fatal(err)
	}

	var response RPCResultResponse
	err = socket.ReadJSON(&response)
	if err != nil {
		t.Fatal(err)
	}

	connectionID := response.Result.(map[string]any)["connection_id"].(float64)
	return socket, int(connectionID)
}

func testFileUpload[Result any](t *testing.T, url string, fields map[string]string, fileName string, filePath string,
	f func(*testing.T, *httptest.ResponseRecorder, *Result, *Error)) {
	t.Run("POST"+url, func(t *testing.T) {

		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		func() {
			defer func() {
				if err := writer.Close(); err != nil {
					t.Fatal(err)
				}
			}()

			for name, value := range fields {
				if err := writer.WriteField(name, value); err != nil {
					t.Fatal(err)
				}
			}

			source, err := files.Fs.OpenFile(filePath, os.O_RDONLY, 0755)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := source.Close(); err != nil {
					t.Fatal(err)
				}
			}()

			uploadWriter, err := writer.CreateFormFile("file", fileName)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := io.Copy(uploadWriter, source); err != nil {
				t.Fatal(err)
			}
		}()

		request := httptest.NewRequest("POST", url, &buf)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		recorder := httptest.NewRecorder()
		HTTPHandler{}.ServeHTTP(recorder, request)

		var errorResponse ErrorResponse
		if err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse); err != nil {
			t.Fatal(err)
			return
		}
		if errorResponse.Error.Code > 0 {
			f(t, recorder, nil, &errorResponse.Error)
			return
		}

		var result Result
		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
			return
		}
		f(t, recorder, &result, nil)
	})
}
