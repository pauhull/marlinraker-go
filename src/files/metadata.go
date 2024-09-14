package files

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Thumbnail struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Size         int    `json:"size"`
	RelativePath string `json:"relative_path"`
}

type Metadata struct {
	FileName            string      `json:"filename"`
	Size                int64       `json:"size"`
	Modified            float64     `json:"modified"`
	PrintStartTime      float64     `json:"print_start_time,omitempty"`
	JobID               string      `json:"job_id,omitempty"`
	Slicer              string      `json:"slicer,omitempty"`
	SlicerVersion       string      `json:"slicer_version,omitempty"`
	LayerHeight         float64     `json:"layer_height,omitempty"`
	FirstLayerHeight    float64     `json:"first_layer_height,omitempty"`
	ObjectHeight        float64     `json:"object_height,omitempty"`
	FilamentTotal       float64     `json:"filament_total,omitempty"`
	EstimatedTime       float64     `json:"estimated_time,omitempty"`
	Thumbnails          []Thumbnail `json:"thumbnails,omitempty"`
	FirstLayerBedTemp   float64     `json:"first_layer_bed_temp,omitempty"`
	FirstLayerExtrTemp  float64     `json:"first_layer_extr_temp,omitempty"`
	GcodeStartByte      int64       `json:"gcode_start_byte,omitempty"`
	GcodeEndByte        int64       `json:"gcode_end_byte,omitempty"`
	NozzleDiameter      float64     `json:"nozzle_diameter,omitempty"`
	FilamentName        string      `json:"filament_name,omitempty"`
	FilamentType        string      `json:"filament_type,omitempty"`
	FilamentWeightTotal float64     `json:"filament_weight_total,omitempty"`
}

func RemoveUnusedMetadata() error {
	gcodeRootPath := filepath.Join(DataDir, "gcodes")
	used := make([]string, 0)

	err := afero.Walk(Fs, gcodeRootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if ext == ".meta" || ext == ".png" {
			return nil
		}

		relPath, err := filepath.Rel(gcodeRootPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		meta, err := LoadMetadata(relPath)
		if err == nil {
			metaDiskPath := getMetaDiskPath(relPath)
			relMetaPath, err := filepath.Rel(gcodeRootPath, metaDiskPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			used = append(used, relMetaPath)
			for _, thumbnail := range meta.Thumbnails {
				used = append(used, filepath.Join(filepath.Dir(relPath), thumbnail.RelativePath))
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk through gcodes directory: %w", err)
	}

	var deleted int
	err = afero.Walk(Fs, gcodeRootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(gcodeRootPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		ext, dir := filepath.Ext(path), filepath.Dir(path)
		if ext == ".meta" {
			if !lo.Contains(used, relPath) {
				if err := Fs.Remove(path); err != nil {
					return fmt.Errorf("failed to remove metadata %q: %w", path, err)
				}
				deleted++
			}
		} else if ext == ".png" && filepath.Base(dir) == ".thumbs" {
			if !lo.Contains(used, relPath) {
				if err := Fs.Remove(path); err != nil {
					return fmt.Errorf("failed to remove thumbnail %q: %w", path, err)
				}
				if isEmpty, err := afero.IsEmpty(Fs, dir); err != nil {
					return fmt.Errorf("failed to check if directory is empty %q: %w", dir, err)
				} else if isEmpty {
					if err := Fs.Remove(dir); err != nil {
						return fmt.Errorf("failed to remove directory %q: %w", dir, err)
					}
				}
				deleted++
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to remove unused metadata: %w", err)
	}

	return nil
}

func HasMetadata(fileName string) bool {
	diskPath := getMetaDiskPath(fileName)
	_, err := Fs.Stat(diskPath)
	return err == nil
}

func LoadOrScanMetadata(fileName string) (*Metadata, error) {

	metadata, err := LoadMetadata(fileName)
	if err == nil {
		log.Debugf("Loaded metadata for %q from disk", fileName)
		return metadata, nil
	}

	metadata, err = ScanMetadata(fileName)
	if err != nil {
		return nil, err
	}
	if err := StoreMetadata(metadata); err != nil {
		return nil, err
	}
	log.Debugf("Scanned and stored metadata for %q", fileName)
	return metadata, nil
}

func StoreMetadata(metadata *Metadata) error {

	diskPath := getMetaDiskPath(metadata.FileName)

	file, err := Fs.OpenFile(diskPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", diskPath, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Failed to close file %q: %v", diskPath, err)
		}
	}()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}
	return nil
}

func LoadMetadata(fileName string) (*Metadata, error) {

	diskPath := getMetaDiskPath(fileName)

	file, err := Fs.Open(diskPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", diskPath, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Failed to close file %q: %v", diskPath, err)
		}
	}()

	metadata := &Metadata{}
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}
	return metadata, nil
}

func MoveMetadata(src string, dest string) error {

	metadata, err := LoadMetadata(src)
	if err != nil {
		return err
	}

	srcMetaDisk, destMetaDisk := getMetaDiskPath(src), getMetaDiskPath(dest)
	srcDir, destDir := filepath.Dir(srcMetaDisk), filepath.Dir(destMetaDisk)
	thumbDir := filepath.Join(destDir, ".thumbs/")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", thumbDir, err)
	}

	destBaseName := filepath.Base(dest)
	for i, thumbnail := range metadata.Thumbnails {
		destRelPath := fmt.Sprintf(".thumbs/%s-%dx%d.png", destBaseName, thumbnail.Width, thumbnail.Height)
		srcDiskPath, destDiskPath := filepath.Join(srcDir, thumbnail.RelativePath), filepath.Join(destDir, destRelPath)
		if err := Fs.Rename(srcDiskPath, destDiskPath); err != nil {
			return fmt.Errorf("failed to move thumbnail %q: %w", srcDiskPath, err)
		}
		thumbnail.RelativePath = destRelPath
		metadata.Thumbnails[i] = thumbnail
	}

	srcThumbDir := filepath.Join(srcDir, ".thumbs/")
	if isEmpty, err := afero.IsEmpty(Fs, srcThumbDir); err != nil {
		return fmt.Errorf("failed to check if directory is empty %q: %w", srcThumbDir, err)
	} else if isEmpty {
		if err := Fs.Remove(srcThumbDir); err != nil {
			return fmt.Errorf("failed to remove directory %q: %w", srcThumbDir, err)
		}
	}

	metadata.FileName = dest
	if err = StoreMetadata(metadata); err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}

	if err = Fs.Remove(srcMetaDisk); err != nil {
		return fmt.Errorf("failed to remove metadata: %w", err)
	}
	return nil
}

func DirectoryRenamed(dir string) error {

	gcodeRootPath := filepath.Join(DataDir, "gcodes")
	diskPath := filepath.Join(gcodeRootPath, dir)
	err := afero.Walk(Fs, diskPath, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			//nolint:wrapcheck
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(gcodeRootPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		metadata, err := LoadMetadata(relPath)
		if err == nil {
			metadata.FileName = relPath
			_ = StoreMetadata(metadata)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk through directory %q: %w", diskPath, err)
	}
	return nil
}

func RemoveMetadata(fileName string) error {

	metadata, err := LoadMetadata(fileName)
	if err != nil {
		return err
	}

	dir := filepath.Join(DataDir, "gcodes", filepath.Dir(fileName))
	for _, thumbnail := range metadata.Thumbnails {
		diskPath := filepath.Join(dir, thumbnail.RelativePath)
		if err := Fs.Remove(diskPath); err != nil {
			return fmt.Errorf("failed to remove thumbnail %q: %w", diskPath, err)
		}
	}

	metaDiskPath := getMetaDiskPath(fileName)
	if err = Fs.Remove(metaDiskPath); err != nil {
		return fmt.Errorf("failed to remove metadata %q: %w", metaDiskPath, err)
	}
	return nil
}

func getMetaDiskPath(fileName string) string {
	gcodeDiskPath := filepath.Join(DataDir, "gcodes", fileName)
	dir, base := filepath.Dir(gcodeDiskPath), filepath.Base(gcodeDiskPath)
	metaDiskPath := filepath.Join(dir, "."+base+".meta")
	return metaDiskPath
}
