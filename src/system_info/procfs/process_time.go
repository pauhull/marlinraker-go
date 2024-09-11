package procfs

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func getProcessTimeImpl(procPidStatPath string, clk int32) (float32, error) {
	statBytes, err := os.ReadFile(procPidStatPath)
	if err != nil {
		return 0, err
	}

	parts := strings.Split(strings.TrimSpace(string(statBytes)), " ")
	if len(parts) < 15 {
		return 0, errors.New("malformed " + procPidStatPath)
	}

	utime, err := strconv.ParseInt(parts[13], 10, 64)
	if err != nil {
		return 0, err
	}

	stime, err := strconv.ParseInt(parts[14], 10, 64)
	if err != nil {
		return 0, err
	}

	return float32(utime+stime) / float32(clk), nil
}
