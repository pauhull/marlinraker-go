package procfs

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type CPUInfo struct {
	CPUCount     int    `json:"cpu_count"`
	Bits         string `json:"bits"`
	Processor    string `json:"processor"`
	CPUDesc      string `json:"cpu_desc"`
	SerialNumber string `json:"serial_number"`
	HardwareDesc string `json:"hardware_desc"`
	Model        string `json:"model"`
	TotalMemory  int64  `json:"total_memory"`
	MemoryUnits  string `json:"memory_units"`
}

var (
	serialRegex    = regexp.MustCompile(`(?m)^Serial\s*: (.+)$`)
	hardwareRegex  = regexp.MustCompile(`(?m)^Hardware\s*: (.+)$`)
	modelRegex     = regexp.MustCompile(`(?m)^Model\s*: (.+)$`)
	modelNameRegex = regexp.MustCompile(`(?m)^model name\s*: (.+)$`)
	processorRegex = regexp.MustCompile(`(?m)^processor\s*: [0-9]+$`)
)

func getCPUInfoImpl(cpuInfoPath string, memInfoPath string) (*CPUInfo, error) {

	info := &CPUInfo{}

	cpuInfoBytes, err := os.ReadFile(cpuInfoPath)
	if err != nil {
		return info, fmt.Errorf("failed to read cpu info: %w", err)
	}
	cpuInfo := string(cpuInfoBytes)

	if matches := processorRegex.FindAllString(cpuInfo, -1); matches != nil {
		info.CPUCount = len(matches)
	} else {
		info.CPUCount = 1
	}

	if out, err := exec.Command("getconf", "LONG_BIT").Output(); err == nil {
		info.Bits = fmt.Sprintf("%sbit", strings.TrimSpace(string(out)))
	} else {
		return info, fmt.Errorf("failed to run 'getconf LONG_BIT': %w", err)
	}

	if out, err := exec.Command("uname", "-m").Output(); err == nil {
		info.Processor = strings.TrimSpace(string(out))
	} else {
		return info, fmt.Errorf("failed to run 'uname -m': %w", err)
	}

	if match := modelNameRegex.FindStringSubmatch(cpuInfo); match != nil {
		info.CPUDesc = match[1]
	}

	if match := serialRegex.FindStringSubmatch(cpuInfo); match != nil {
		info.SerialNumber = match[1]
	}

	if match := hardwareRegex.FindStringSubmatch(cpuInfo); match != nil {
		info.HardwareDesc = match[1]
	}

	if match := modelRegex.FindStringSubmatch(cpuInfo); match != nil {
		info.Model = match[1]
	}

	if info.TotalMemory, info.MemoryUnits, err = getTotalMemImpl(memInfoPath); err != nil {
		info.TotalMemory = 0
		info.MemoryUnits = "B"
	}

	return info, nil
}

type CPUTimes = map[string][]int64

var cpuRegex = regexp.MustCompile(`(?m)^(cpu[0-9]*)\s*(.*)$`)

func getCPUTimesImpl(procStatPath string) (CPUTimes, error) {

	times := make(CPUTimes)

	procStatBytes, err := os.ReadFile(procStatPath)
	if err != nil {
		return times, fmt.Errorf("failed reading %q: %w", procStatPath, err)
	}
	procStats := string(procStatBytes)

	if matches := cpuRegex.FindAllStringSubmatch(procStats, -1); matches != nil {

		for _, match := range matches {
			cpu := match[1]
			parts := strings.Split(match[2], " ")
			_times := make([]int64, len(parts))

			for index, part := range parts {
				time, err := strconv.ParseInt(part, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse cpu time: %w", err)
				}
				_times[index] = time
			}

			times[cpu] = _times
		}
	}

	return times, nil
}

type CPUUsage = map[string]float64

func GetCPUUsage(lastTimes CPUTimes, currentTimes CPUTimes) CPUUsage {

	usage := make(CPUUsage)

	for cpu, current := range currentTimes {
		var (
			last = lastTimes[cpu]
			sum  int64
			idle int64
		)
		for i := 0; i < len(current); i++ {
			sum += current[i] - last[i]
			if i == 3 {
				idle = current[i] - last[i]
			}
		}
		usage[cpu] = 100.0 - math.Floor(float64(idle)/float64(sum)*1000.0)/10.0
	}

	return usage
}

func getCPUTempImpl(thermalZonePath string) (float64, error) {

	tempBytes, err := os.ReadFile(thermalZonePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %q: %w", thermalZonePath, err)
	}

	temp, err := strconv.ParseInt(strings.TrimSpace(string(tempBytes)), 10, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to parse cpu temp: %w", err)
	}

	return math.Round(float64(temp)/10.0) / 100.0, nil
}
