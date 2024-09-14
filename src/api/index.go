package api

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"marlinraker/src/constants"
	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/system_info"
)

var (
	//go:embed index.html
	indexHTML     string
	indexTemplate = template.Must(template.New("index.html").Parse(indexHTML))
)

type templateData struct {
	Hostname       string
	RequestIP      string
	Version        string
	WebsocketCount int
	State          marlinraker.KlippyState
	Arch           string
	Os             string
	CPU            string
}

func handleIndex(writer http.ResponseWriter, request *http.Request) error {

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	info, err := system_info.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}

	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(200)

	err = indexTemplate.Execute(writer, templateData{
		Hostname:       hostname,
		RequestIP:      getAddress(request),
		Version:        constants.Version,
		WebsocketCount: len(connections.GetConnections()),
		State:          marlinraker.State,
		Arch:           info.CPUInfo.Processor,
		Os:             info.Distribution.Name,
		CPU:            info.CPUInfo.CPUDesc,
	})
	if err != nil {
		return fmt.Errorf("failed to substitute template: %w", err)
	}
	return nil
}

func getAddress(request *http.Request) string {
	address := request.Header.Get("X-Real-Ip")
	if address == "" {
		address = request.Header.Get("X-Forwarded-For")
	}
	if address == "" {
		address = request.RemoteAddr
	}
	return address
}
