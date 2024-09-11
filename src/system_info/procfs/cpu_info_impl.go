//go:build linux

package procfs

func GetCpuInfo() (*CpuInfo, error) {
	return getCpuInfoImpl("/proc/cpuinfo", "/proc/meminfo")
}

func GetCpuTimes() (CpuTimes, error) {
	return getCpuTimesImpl("/proc/stat")
}

func GetCpuTemp() (float64, error) {
	return getCpuTempImpl("/sys/class/thermal/thermal_zone0/temp")
}
