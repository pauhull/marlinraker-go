//go:build linux

package procfs

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetProcessTime() (float32, error) {
	clkBytes, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		return 0, err
	}
	clk, err := strconv.ParseInt(strings.TrimSpace(string(clkBytes)), 10, 32)

	pid := os.Getpid()
	return getProcessTimeImpl("/proc/"+strconv.Itoa(pid)+"/stat", int32(clk))
}
