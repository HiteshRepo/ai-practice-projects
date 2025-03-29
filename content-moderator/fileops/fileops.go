package fileops

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(fileloc string) (string, error) {
	if fileloc == "" {
		return "", fmt.Errorf("filepath cannot be empty")
	}

	err := CheckFile(fileloc)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(fileloc)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return string(data), nil
}

func CheckFile(fileloc string) error {
	absPath, err := filepath.Abs(fileloc)
	if err != nil {
		return fmt.Errorf("invalid path format: %w", err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", absPath)
		}

		return fmt.Errorf("unable to access file: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", absPath)
	}

	return nil
}
