package files

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"marlinraker/src/api/notification"
	"marlinraker/src/util"
)

type FileMeta struct {
	Modified    float64 `json:"modified"`
	Size        int64   `json:"size"`
	Permissions string  `json:"permissions"`
	FileName    string  `json:"filename"`
}

type ExtendedFileMeta struct {
	*Metadata
	Permissions string `json:"permissions"`
}

type FileUploadAction struct {
	Item         ActionItem `json:"item"`
	Action       string     `json:"action"`
	PrintStarted *bool      `json:"print_started,omitempty"`
}

type FileDeleteAction struct {
	Item   ActionItem `json:"item"`
	Action string     `json:"action"`
}

func Upload(rootName, path, checksum string, header *multipart.FileHeader) (FileUploadAction, error) {

	fileName := header.Filename
	sourceFile, err := header.Open()
	if err != nil {
		return FileUploadAction{}, fmt.Errorf("failed to open file %q: %w", fileName, err)
	}

	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Errorf("Failed to close file %q: %v", fileName, err)
		}
	}()

	root, err := getRootByName(rootName)
	if err != nil {
		return FileUploadAction{}, err
	}

	if !strings.Contains(root.Permissions, "w") {
		return FileUploadAction{}, errors.New("no write permissions")
	}

	destDirPath := filepath.Join(root.Path, path)
	destPath := filepath.Join(destDirPath, fileName)

	if err := Fs.MkdirAll(destDirPath, 0755); err != nil {
		return FileUploadAction{}, fmt.Errorf("failed to create directory %q: %w", destDirPath, err)
	}

	_, err = Fs.Stat(destPath)
	existedBefore := err == nil

	destFile, err := Fs.OpenFile(destPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return FileUploadAction{}, fmt.Errorf("failed to open file %q: %w", destPath, err)
	}

	defer func() {
		if err := destFile.Close(); err != nil {
			log.Errorf("Failed to close file %q: %v", destPath, err)
		}
	}()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return FileUploadAction{}, fmt.Errorf("failed to copy file %q to %q: %w", fileName, destPath, err)
	}

	if checksum != "" {
		hash, err := calculateChecksum(destPath)
		if err != nil {
			return FileUploadAction{}, err
		}
		if checksum != hash {
			return FileUploadAction{}, util.NewErrorf(422, "checksums do not match (got %s, expected: %s)", checksum, hash)
		}
	}

	stat, err := Fs.Stat(destPath)
	if err != nil {
		return FileUploadAction{}, fmt.Errorf("failed to stat file %q: %w", destPath, err)
	}

	actionName := "create_file"
	if existedBefore {
		actionName = "modify_file"
	}

	action := FileUploadAction{
		Item: ActionItem{
			Path:        filepath.Join(path, fileName),
			Root:        root.Name,
			Modified:    float64(stat.ModTime().UnixMilli()) / 1000.0,
			Size:        stat.Size(),
			Permissions: root.Permissions,
		},
		Action: actionName,
	}

	err = notification.Publish(notification.New("notify_filelist_changed", []any{action}))
	if err != nil {
		return FileUploadAction{}, fmt.Errorf("failed to publish notification: %w", err)
	}
	return action, nil
}

func DeleteFile(path string) (FileDeleteAction, error) {

	idx := strings.IndexByte(path, '/')
	if idx == -1 {
		return FileDeleteAction{}, errors.New("invalid filepath")
	}
	rootName, fileName := path[:idx], path[idx+1:]

	root, err := getRootByName(rootName)
	if err != nil {
		return FileDeleteAction{}, err
	}

	if !strings.Contains(root.Permissions, "w") {
		return FileDeleteAction{}, errors.New("no write permissions")
	}

	diskPath := filepath.Join(DataDir, rootName, fileName)
	stat, err := Fs.Stat(diskPath)
	if err != nil {
		return FileDeleteAction{}, fmt.Errorf("failed to stat file %q: %w", diskPath, err)
	}
	if stat.IsDir() {
		return FileDeleteAction{}, util.NewErrorf(400, "%s is a directory", diskPath)
	}
	if err := Fs.Remove(diskPath); err != nil {
		return FileDeleteAction{}, fmt.Errorf("failed to delete file %q: %w", diskPath, err)
	}

	if HasMetadata(fileName) {
		if err := RemoveMetadata(fileName); err != nil {
			return FileDeleteAction{}, err
		}
	}

	action := FileDeleteAction{
		Item: ActionItem{
			Path:        fileName,
			Root:        rootName,
			Size:        0,
			Modified:    0,
			Permissions: "",
		},
		Action: "delete_file",
	}

	err = notification.Publish(notification.New("notify_filelist_changed", []any{action}))
	if err != nil {
		return FileDeleteAction{}, fmt.Errorf("failed to publish notification: %w", err)
	}
	return action, nil
}

func calculateChecksum(filePath string) (string, error) {

	file, err := Fs.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %w", filePath, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Failed to close file %q: %v", filePath, err)
		}
	}()

	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	checksum := fmt.Sprintf("%x", digest.Sum(nil))
	return checksum, nil
}
