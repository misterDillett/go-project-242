package code

import (
	"fmt"
	"os"
	"strings"
)

func GetSize(path string, includeHidden bool) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	var totalSize int64
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}

		if !entryInfo.IsDir() {
			totalSize += entryInfo.Size()
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
	case 5:
		return fmt.Sprintf("%.1fEB", float64(size)/float64(div))
	default:
		return fmt.Sprintf("%.1fEB", float64(size)/float64(div))
	}
}