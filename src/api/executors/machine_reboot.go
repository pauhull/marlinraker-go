package executors

import (
	log "github.com/sirupsen/logrus"
	"marlinraker/src/marlinraker/connections"
	"net/http"
	"os/exec"
)

type MachineRebootResult string

func MachineReboot(*connections.Connection, *http.Request, Params) (any, error) {
	go func() {
		connections.TerminateAllConnections()
		if out, err := exec.Command("systemctl", "reboot").CombinedOutput(); err != nil {
			log.Errorf("Failed to reboot machine: %v", err)
			log.Errorln(string(out))
		}
	}()
	return "ok", nil
}
