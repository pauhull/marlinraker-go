package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"marlinraker/src/api/executors"
	"marlinraker/src/files"
	"marlinraker/src/marlinraker"
	"marlinraker/src/marlinraker/connections"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

type Executor func(*connections.Connection, *http.Request, executors.Params) (any, error)

var socketExecutors = map[string]Executor{
	"access.oneshot_token":          executors.AccessOneshotToken,
	"machine.proc_stats":            executors.MachineProcStats,
	"machine.reboot":                executors.MachineReboot,
	"machine.services.restart":      executors.MachineServicesRestart,
	"machine.services.start":        executors.MachineServicesStart,
	"machine.services.stop":         executors.MachineServicesStop,
	"machine.shutdown":              executors.MachineShutdown,
	"machine.system_info":           executors.MachineSystemInfo,
	"printer.emergency_stop":        executors.PrinterEmergencyStop,
	"printer.firmware_restart":      executors.PrinterFirmwareRestart,
	"printer.gcode.help":            executors.PrinterGcodeHelp,
	"printer.gcode.script":          executors.PrinterGcodeScript,
	"printer.info":                  executors.PrinterInfo,
	"printer.objects.list":          executors.PrinterObjectsList,
	"printer.objects.query":         executors.PrinterObjectsQuerySocket,
	"printer.objects.subscribe":     executors.PrinterObjectsSubscribeSocket,
	"printer.print.cancel":          executors.PrinterPrintCancel,
	"printer.print.pause":           executors.PrinterPrintPause,
	"printer.print.resume":          executors.PrinterPrintResume,
	"printer.print.start":           executors.PrinterPrintStart,
	"printer.restart":               executors.PrinterRestart,
	"server.config":                 executors.ServerConfig,
	"server.connection.identify":    executors.ServerConnectionIdentify,
	"server.database.delete_item":   executors.ServerDatabaseDeleteItem,
	"server.database.get_item":      executors.ServerDatabaseGetItem,
	"server.database.list":          executors.ServerDatabaseList,
	"server.database.post_item":     executors.ServerDatabasePostItem,
	"server.files.delete_directory": executors.ServerFilesDeleteDirectory,
	"server.files.delete_file":      executors.ServerFilesDeleteFile,
	"server.files.get_directory":    executors.ServerFilesGetDirectory,
	"server.files.list":             executors.ServerFilesList,
	"server.files.metadata":         executors.ServerFilesMetadata,
	"server.files.metascan":         executors.ServerFilesMetadata,
	"server.files.move":             executors.ServerFilesMove,
	"server.files.post_directory":   executors.ServerFilesPostDirectory,
	"server.files.roots":            executors.ServerFilesRoots,
	"server.files.thumbnails":       executors.ServerFilesThumbnails,
	"server.files.zip":              executors.ServerFilesZip,
	"server.gcode_store":            executors.ServerGcodeStore,
	"server.history.list":           executors.ServerHistoryList,
	"server.info":                   executors.ServerInfo,
	"server.logs.rollover":          executors.ServerLogsRollover,
	"server.restart":                executors.ServerRestart,
	"server.temperature_store":      executors.ServerTemperatureStore,
	"server.webcams.list":           executors.ServerWebcamsList,
}

var httpExecutors = map[string]map[string]Executor{
	"GET": {
		"/access/oneshot_token":     executors.AccessOneshotToken,
		"/machine/proc_stats":       executors.MachineProcStats,
		"/machine/system_info":      executors.MachineSystemInfo,
		"/printer/gcode/help":       executors.PrinterGcodeHelp,
		"/printer/info":             executors.PrinterInfo,
		"/printer/objects/list":     executors.PrinterObjectsList,
		"/printer/objects/query":    executors.PrinterObjectsQueryHTTP,
		"/server/config":            executors.ServerConfig,
		"/server/database/item":     executors.ServerDatabaseGetItem,
		"/server/database/list":     executors.ServerDatabaseList,
		"/server/files/directory":   executors.ServerFilesGetDirectory,
		"/server/files/list":        executors.ServerFilesList,
		"/server/files/metadata":    executors.ServerFilesMetadata,
		"/server/files/metascan":    executors.ServerFilesMetascan,
		"/server/files/roots":       executors.ServerFilesRoots,
		"/server/files/thumbnails":  executors.ServerFilesThumbnails,
		"/server/gcode_store":       executors.ServerGcodeStore,
		"/server/history/list":      executors.ServerHistoryList,
		"/server/info":              executors.ServerInfo,
		"/server/temperature_store": executors.ServerTemperatureStore,
	},
	"POST": {
		"/machine/reboot":            executors.MachineReboot,
		"/machine/services/restart":  executors.MachineServicesRestart,
		"/machine/services/start":    executors.MachineServicesStart,
		"/machine/services/stop":     executors.MachineServicesStop,
		"/machine/shutdown":          executors.MachineShutdown,
		"/printer/emergency_stop":    executors.PrinterEmergencyStop,
		"/printer/firmware_restart":  executors.PrinterFirmwareRestart,
		"/printer/gcode/script":      executors.PrinterGcodeScript,
		"/printer/objects/subscribe": executors.PrinterObjectsSubscribeHTTP,
		"/printer/print/cancel":      executors.PrinterPrintCancel,
		"/printer/print/pause":       executors.PrinterPrintPause,
		"/printer/print/resume":      executors.PrinterPrintResume,
		"/printer/print/start":       executors.PrinterPrintStart,
		"/printer/restart":           executors.PrinterRestart,
		"/server/database/item":      executors.ServerDatabasePostItem,
		"/server/files/directory":    executors.ServerFilesPostDirectory,
		"/server/files/move":         executors.ServerFilesMove,
		"/server/files/upload":       executors.ServerFilesUpload,
		"/server/files/zip":          executors.ServerFilesZip,
		"/server/logs/rollover":      executors.ServerLogsRollover,
		"/server/restart":            executors.ServerRestart,
	},
	"DELETE": {
		"/server/database/item":   executors.ServerDatabaseDeleteItem,
		"/server/files/directory": executors.ServerFilesDeleteDirectory,
	},
}

type HTTPHandler struct{}

func (HTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log.WithField("method", request.Method).Debugln(request.URL.String())

	requestPath := strings.TrimRight(request.URL.Path, "/")
	err := handlePath(writer, request, requestPath)
	if err != nil {
		log.Errorf("Error while handling %s %s: %v", request.Method, requestPath, err)
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
	}
}

func handlePath(writer http.ResponseWriter, request *http.Request, requestPath string) error {
	switch {

	case requestPath == "":
		return handleIndex(writer, request)

	case requestPath == "/websocket":
		return handleSocket(writer, request)

	case marlinraker.Config.Misc.OctoprintCompat && strings.HasPrefix(requestPath, "/api/"):
		return handleOctoPrint(writer, request)

	case isFilePath(requestPath) && (request.Method == "GET" || request.Method == "DELETE"):
		if request.Method == "GET" {
			return handleFileDownload(writer, request)
		} else {
			return handleFileDelete(writer, request)
		}

	default:
		return handleHTTP(writer, request)
	}
}

func StartServer() {
	address := fmt.Sprintf("%s:%d", marlinraker.Config.Web.BindAddress, marlinraker.Config.Web.Port)
	log.Printf("Listening on %s", address)

	server := &http.Server{
		Addr:              address,
		Handler:           cors.AllowAll().Handler(HTTPHandler{}),
		ReadHeaderTimeout: 0, // we might have long-running requests
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func isFilePath(path string) bool {
	for _, root := range files.FileRoots {
		if strings.HasPrefix(path, "/server/files/"+root.Name) {
			return true
		}
	}
	return lo.Contains([]string{"/server/files/klippy.log", "/server/files/moonraker.log", "/server/files/marlinraker.log"}, path)
}
