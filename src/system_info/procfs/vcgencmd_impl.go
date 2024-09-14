//go:build linux

package procfs

import (
	"fmt"
	"os/exec"
)

func GetThrottledState() (ThrottledState, error) {
	throttledBytes, err := exec.Command("vcgencmd", "get_throttled").Output()
	if err != nil {
		return ThrottledState{0, []string{}}, fmt.Errorf("failed to run vcgencmd: %w", err)
	}
	return getThrottledStateImpl(throttledBytes)
}
