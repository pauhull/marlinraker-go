package procfs

import (
	"fmt"
	"github.com/samber/lo"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
)

type IpAddress struct {
	Family      string `json:"family"`
	Address     string `json:"address"`
	IsLinkLocal bool   `json:"is_link_local"`
}

type Network struct {
	MacAddress  string      `json:"mac_address"`
	IpAddresses []IpAddress `json:"ip_addresses"`
}

func GetNetwork() (map[string]Network, error) {

	networks := make(map[string]Network)

	ifaces, err := net.Interfaces()
	if err != nil {
		return networks, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range ifaces {

		addresses, err := iface.Addrs()
		if err != nil {
			return networks, fmt.Errorf("failed to get addresses for %q: %w", iface.Name, err)
		}

		networks[iface.Name] = Network{
			MacAddress: iface.HardwareAddr.String(),
			IpAddresses: lo.Map(addresses, func(address net.Addr, _ int) IpAddress {

				if ipv4 := address.(*net.IPNet).IP.To4(); ipv4 != nil {
					return IpAddress{
						Address:     ipv4.String(),
						Family:      "ipv4",
						IsLinkLocal: ipv4.IsLinkLocalUnicast() || ipv4.IsLinkLocalMulticast(),
					}
				}

				if ipv6 := address.(*net.IPNet).IP.To16(); ipv6 != nil {
					return IpAddress{
						Address:     ipv6.String(),
						Family:      "ipv6",
						IsLinkLocal: ipv6.IsLinkLocalUnicast() || ipv6.IsLinkLocalMulticast(),
					}
				}

				return IpAddress{}
			}),
		}
	}

	return networks, nil
}

type ifaceStats struct {
	RxBytes   int64   `json:"rx_bytes"`
	TxBytes   int64   `json:"tx_bytes"`
	Bandwidth float64 `json:"bandwidth"`
}

type NetworkStats map[string]ifaceStats

type TimedNetworkStats struct {
	Stats NetworkStats
	Time  float64
}

var (
	ifaceRegex      = regexp.MustCompile(`(?m)^\s*(.*):\s*([0-9 ]+?)\s*$`)
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

func getNetworkStatsImpl(lastStats *TimedNetworkStats, now float64, procNetDevPath string) (*TimedNetworkStats, error) {

	procNetDevBytes, err := os.ReadFile(procNetDevPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", procNetDevPath, err)
	}
	procNetDev := string(procNetDevBytes)

	stats := make(NetworkStats)
	if matches := ifaceRegex.FindAllStringSubmatch(procNetDev, -1); matches != nil {
		for _, match := range matches {

			iface, data := match[1], match[2]
			nums := whitespaceRegex.Split(data, -1)
			rxBytes, _ := strconv.ParseInt(nums[0], 10, 0)
			txBytes, _ := strconv.ParseInt(nums[8], 10, 0)

			bandwidth := 0.0
			if lastStats != nil {
				if last, exists := lastStats.Stats[iface]; exists {
					delta := rxBytes + txBytes - last.RxBytes - last.TxBytes
					bandwidth = math.Floor(float64(delta) / (now - lastStats.Time) * 100.0)
				}
			}

			stats[iface] = ifaceStats{
				RxBytes:   rxBytes,
				TxBytes:   txBytes,
				Bandwidth: bandwidth,
			}
		}
	}

	return &TimedNetworkStats{Stats: stats, Time: now}, nil
}
