package constants

import (
	"strconv"
	"strings"
)

var (
	APIVersionString = "0.0.0"
	suffix           = "dev"

	Version = func() string {
		if suffix != "" {
			return APIVersionString + "-" + suffix
		}
		return APIVersionString
	}()

	APIVersion = func() [3]int {
		parts := strings.Split(APIVersionString, ".")
		major, _ := strconv.Atoi(parts[0])
		minor, _ := strconv.Atoi(parts[1])
		patch, _ := strconv.Atoi(parts[2])
		return [3]int{major, minor, patch}
	}()
)
