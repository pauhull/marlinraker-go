package config

import (
	"github.com/BurntSushi/toml"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"marlinraker/src/files"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Web struct {
	BindAddress string   `toml:"bind_address"`
	Port        int      `toml:"port"`
	CorsDomains []string `toml:"cors_domains"`
}

type Serial struct {
	Port                  string      `toml:"port"`
	BaudRate              interface{} `toml:"baud_rate"`
	MaxConnectionAttempts int         `toml:"max_connection_attempts"`
	ConnectionTimeout     int         `toml:"connection_timeout"`
}

type Misc struct {
	OctoprintCompat bool     `toml:"octoprint_compat"`
	ExtendedLogs    bool     `toml:"extended_logs"`
	AllowedServices []string `toml:"allowed_services"`
}

type Heater struct {
	MinTemp int `toml:"min_temp"`
	MaxTemp int `toml:"max_temp"`
}

type Extruder struct {
	Heater
	MinExtrudeTemp   int     `toml:"min_extrude_temp"`
	FilamentDiameter float32 `toml:"filament_diameter"`
}

type HeaterBed struct {
	Heater
}

type Gcode struct {
	SendM73        bool `toml:"send_m73"`
	ReportVelocity bool `toml:"report_velocity"`
}

type Printer struct {
	BedMesh     bool      `toml:"bed_mesh"`
	AxisMinimum [3]int    `toml:"axis_minimum"`
	AxisMaximum [3]int    `toml:"axis_maximum"`
	Extruder    Extruder  `toml:"extruder"`
	HeaterBed   HeaterBed `toml:"heater_bed"`
	Gcode       Gcode     `toml:"gcode"`
}

type Macro struct {
	RenameExisting string         `toml:"rename_existing"`
	Variables      map[string]any `toml:"variables"`
	Gcode          string         `toml:"gcode"`
}

type Config struct {
	Web     Web              `toml:"web"`
	Serial  Serial           `toml:"serial"`
	Misc    Misc             `toml:"misc"`
	Printer Printer          `toml:"printer"`
	Macros  map[string]Macro `toml:"macros"`
}

var includeRegex = regexp.MustCompile(`(?mi)^#include +(\S+).*$`)

func LoadConfig(path string) (*Config, error) {
	if _, err := files.Fs.Stat(path); err != nil {
		return nil, err
	}
	contents, err := resolve(path, []string{})
	if err != nil {
		return nil, err
	}
	return parseConfig(contents)
}

func resolve(currentPath string, resolvedSoFar []string) (string, error) {

	if !filepath.IsAbs(currentPath) {
		currentPath, _ = filepath.Abs(currentPath)
	}

	stat, err := files.Fs.Stat(currentPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warnf("File or directory %q not found", currentPath)
			return "", nil
		}
		return "", err
	}

	if stat.IsDir() {
		dirFiles, err := afero.ReadDir(files.Fs, currentPath)
		if err != nil {
			return "", err
		}
		contents := make([]string, 0)
		for _, file := range dirFiles {
			if file.IsDir() {
				continue
			}
			content, err := resolve(filepath.Join(currentPath, file.Name()), resolvedSoFar)
			if err != nil {
				return "", err
			}
			contents = append(contents, content)
		}
		return strings.Join(contents, "\n"), nil

	} else {

		contentBytes, err := afero.ReadFile(files.Fs, currentPath)
		if err != nil {
			return "", err
		}
		content := string(contentBytes)
		resolvedSoFar = append(resolvedSoFar, currentPath)

		matches := includeRegex.FindAllStringIndex(content, -1)
		for i := len(matches) - 1; i >= 0; i-- {
			start, end := matches[i][0], matches[i][1]
			str := content[start:end]
			filename := includeRegex.FindStringSubmatch(str)[1]
			nextPath, err := filepath.Abs(filepath.Join(filepath.Dir(currentPath), filename))
			if err != nil {
				return "", err
			}

			if lo.Contains(resolvedSoFar, nextPath) {
				log.Warnf("Cannot resolve cyclic dependency %q in %q", filename, currentPath)
				continue
			}

			result, err := resolve(nextPath, resolvedSoFar)
			if err != nil {
				return "", err
			}

			content = content[0:start] + result + content[end:]
		}

		return content, nil
	}
}

func DefaultConfig() *Config {
	return &Config{
		Web: Web{
			Port:        7125,
			CorsDomains: []string{},
		},
		Serial: Serial{
			Port:                  "auto",
			BaudRate:              "auto",
			MaxConnectionAttempts: 5,
			ConnectionTimeout:     5000,
		},
		Misc: Misc{
			OctoprintCompat: true,
			ExtendedLogs:    false,
			AllowedServices: []string{"marlinraker", "crowsnest", "MoonCord", "moonraker-telegram-bot", "KlipperScreen", "sonar", "webcamd"},
		},
		Printer: Printer{
			BedMesh:     false,
			AxisMinimum: [3]int{0, 0, 0},
			AxisMaximum: [3]int{220, 220, 240},
			Extruder: Extruder{
				Heater: Heater{
					MinTemp: 0,
					MaxTemp: 250,
				},
				MinExtrudeTemp:   180,
				FilamentDiameter: 1.75,
			},
			HeaterBed: HeaterBed{
				Heater: Heater{
					MinTemp: 0,
					MaxTemp: 100,
				},
			},
			Gcode: Gcode{
				SendM73:        true,
				ReportVelocity: true,
			},
		},
		Macros: map[string]Macro{},
	}
}

func parseConfig(contents string) (*Config, error) {
	config := DefaultConfig()
	_, err := toml.Decode(contents, config)
	return config, err
}
