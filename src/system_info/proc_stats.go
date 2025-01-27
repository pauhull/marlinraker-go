package system_info

import (
	"fmt"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"marlinraker/src/api/notification"
	"marlinraker/src/marlinraker/connections"
	"marlinraker/src/printer_objects"
	"marlinraker/src/system_info/procfs"
	"sync"
	"time"
)

type snapshot struct {
	Time     float64 `json:"time"`
	CpuUsage float64 `json:"cpu_usage"`
	Memory   int64   `json:"memory"`
	MemUnits string  `json:"mem_units"`
}

type ProcStats struct {
	MoonrakerStats       []snapshot            `json:"moonraker_stats"`
	ThrottledState       procfs.ThrottledState `json:"throttled_state"`
	CpuTemp              float64               `json:"cpu_temp"`
	Network              procfs.NetworkStats   `json:"network"`
	SystemCpuUsage       procfs.CpuUsage       `json:"system_cpu_usage"`
	SystemUptime         float64               `json:"system_uptime"`
	WebsocketConnections int                   `json:"websocket_connections"`
}

type ProcStat struct {
	MoonrakerStats       snapshot            `json:"moonraker_stats"`
	CpuTemp              float64             `json:"cpu_temp"`
	Network              procfs.NetworkStats `json:"network"`
	SystemCpuUsage       procfs.CpuUsage     `json:"system_cpu_usage"`
	WebsocketConnections int                 `json:"websocket_connections"`
}

var (
	lastTimes        procfs.CpuTimes
	lastNetworkStats *procfs.TimedNetworkStats
	statsMutex       = &sync.RWMutex{}
	stats            = &ProcStats{}
)

func Run() {
	printer_objects.RegisterObject("system_stats", systemStatsObject{})

	var err error
	lastTimes, err = procfs.GetCpuTimes()
	if err != nil {
		log.Errorf("Failed to get CPU times: %v", err)
		return
	}

	lastNetworkStats, err = procfs.GetNetworkStats(nil)
	if err != nil {
		log.Errorf("Failed to get network stats: %v", err)
		return
	}

	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if err := printer_objects.EmitObject("system_stats"); err != nil {
			log.Errorf("Failed to emit system_stats object: %v", err)
			break
		}
		if err := takeSnapshot(); err != nil {
			log.Errorf("Failed to take system stats snapshot: %v", err)
			break
		}
	}
}

func GetStats() *ProcStats {
	statsMutex.RLock()
	defer statsMutex.RUnlock()
	return stats
}

func takeSnapshot() error {

	statsMutex.Lock()
	defer statsMutex.Unlock()

	now := float64(time.Now().UnixMilli()) / 1000.0

	currentTimes, err := procfs.GetCpuTimes()
	if err != nil {
		return fmt.Errorf("failed to get CPU times: %w", err)
	}

	usedMem, memUnits, err := procfs.GetUsedMem()
	if err != nil {
		return fmt.Errorf("failed to get used memory: %w", err)
	}

	stats.SystemCpuUsage = procfs.GetCpuUsage(lastTimes, currentTimes)
	lastTimes = currentTimes
	avgCpuUsage := stats.SystemCpuUsage["cpu"]

	stats.MoonrakerStats = append(stats.MoonrakerStats, snapshot{
		Time:     now,
		CpuUsage: avgCpuUsage,
		Memory:   usedMem,
		MemUnits: memUnits,
	})

	if len(stats.MoonrakerStats) > 30 {
		stats.MoonrakerStats = stats.MoonrakerStats[1:]
	}

	stats.CpuTemp, err = procfs.GetCpuTemp()
	if err != nil {
		return fmt.Errorf("failed to get CPU temp: %w", err)
	}

	networkStats, err := procfs.GetNetworkStats(lastNetworkStats)
	if err != nil {
		return fmt.Errorf("failed to get network stats: %w", err)
	}
	lastNetworkStats = networkStats
	stats.Network = networkStats.Stats

	stats.SystemUptime, err = procfs.GetUptime()
	if err != nil {
		return fmt.Errorf("failed to get system uptime: %w", err)
	}

	stats.WebsocketConnections = len(connections.GetConnections())

	throttledState, _ := procfs.GetThrottledState()
	if stats.ThrottledState.Bits != throttledState.Bits {
		err = notification.Publish(notification.New("notify_proc_stat_update", []any{throttledState}))
		if err != nil {
			return fmt.Errorf("failed to publish notification: %w", err)
		}
	}
	stats.ThrottledState = throttledState

	notify := notification.New("notify_proc_stat_update", []any{ProcStat{
		MoonrakerStats:       lo.Must(lo.Last(stats.MoonrakerStats)),
		CpuTemp:              stats.CpuTemp,
		Network:              stats.Network,
		SystemCpuUsage:       stats.SystemCpuUsage,
		WebsocketConnections: stats.WebsocketConnections,
	}})

	err = notification.Publish(notify)
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	return nil
}
