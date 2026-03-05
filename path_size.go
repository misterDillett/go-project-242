package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetSize(path string, recursive bool, includeHidden bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	return getDirSize(path, recursive, includeHidden)
}

func getDirSize(path string, recursive bool, includeHidden bool) (int64, error) {
	var totalSize int64

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		info, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		if !info.IsDir() {
			totalSize += info.Size()
		} else if recursive {
			subSize, err := getDirSize(fullPath, recursive, includeHidden)
			if err == nil {
				totalSize += subSize
			}
		}
	}

	return totalSize, nil
}

func FormatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	switch exp {
	case 0:
		return fmt.Sprintf("%.1fKB", float64(size)/float64(unit))
	case 1:
		return fmt.Sprintf("%.1fMB", float64(size)/float64(div))
	case 2:
		return fmt.Sprintf("%.1fGB", float64(size)/float64(div))
	case 3:
		return fmt.Sprintf("%.1fTB", float64(size)/float64(div))
	case 4:
		return fmt.Sprintf("%.1fPB", float64(size)/float64(div))
	default:
		return fmt.Sprintf("%.1fEB", float64(size)/float64(div))
	}
}