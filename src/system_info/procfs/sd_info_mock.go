//go:build !linux

package procfs

func GetSdInfo() (*SdInfo, error) {
	return &SdInfo{
		ManufacturerID:   "",
		Manufacturer:     "Unknown",
		OemID:            "",
		ProductName:      "Unknown",
		ProductRevision:  "",
		SerialNumber:     "",
		ManufacturerDate: "",
		Capacity:         "",
		TotalBytes:       0,
	}, nil
}
