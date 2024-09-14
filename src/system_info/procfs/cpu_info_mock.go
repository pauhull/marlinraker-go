//go:build !linux

package procfs

import (
	"runtime"
)

func GetCPUInfo() (*CPUInfo, error) {
	return &CPUInfo{
		CPUCount:     runtime.NumCPU(),
		Bits:         "",
		Processor:    "Unknown",
		CPUDesc:      "",
		SerialNumber: "Unknown",
		HardwareDesc: "",
		Model:        "Unknown",
		TotalMemory:  0,
		MemoryUnits:  "B",
	}, nil
}

func GetCPUTimes() (CPUTimes, error) {
	return CPUTimes{}, nil
}

func GetCPUTemp() (float64, error) {
	return 0, nil
}
