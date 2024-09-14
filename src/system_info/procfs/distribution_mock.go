//go:build !linux

package procfs

func GetDistribution() (*Distribution, error) {
	return &Distribution{
		Name:         "Unknown",
		ID:           "",
		Version:      "",
		VersionParts: versionParts{},
		Like:         "",
		Codename:     "",
	}, nil
}
