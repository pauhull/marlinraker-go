package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var m114Regex = regexp.MustCompile(`([XYZE]):([+-]?[0-9.]+)`)

func ParseM114(response string) ([4]float64, error) {
	position := [4]float64{}
	if idx := strings.Index(response, "Count"); idx != -1 {
		response = response[:idx]
	}

	for _, match := range m114Regex.FindAllStringSubmatch(response, -1) {
		val, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			return position, fmt.Errorf("failed to parse %s coordinate: %w", match[1], err)
		}
		idx := strings.Index("XYZE", match[1])
		position[idx] = val
	}
	return position, nil
}
