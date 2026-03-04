package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSize_File(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte("hello world")
	err = os.WriteFile(tmpFile.Name(), content, 0644)
	require.NoError(t, err)

	size, err := GetSize(tmpFile.Name(), false)
	require.NoError(t, err)
	assert.Equal(t, int64(len(content)), size)
}

func TestGetSize_Directory_WithoutHidden(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.txt")
	err = os.WriteFile(file1, []byte("hello"), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "file2.txt")
	err = os.WriteFile(file2, []byte("world"), 0644)
	require.NoError(t, err)

	hiddenFile := filepath.Join(tmpDir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("secret"), 0644)
	require.NoError(t, err)

	subdir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subdir, 0755)
	require.NoError(t, err)

	size, err := GetSize(tmpDir, false)
	require.NoError(t, err)
	assert.Equal(t, int64(10), size)
}

func TestGetSize_Directory_WithHidden(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.txt")
	err = os.WriteFile(file1, []byte("hello"), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "file2.txt")
	err = os.WriteFile(file2, []byte("world"), 0644)
	require.NoError(t, err)

	hiddenFile := filepath.Join(tmpDir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("secret"), 0644)
	require.NoError(t, err)

	subdir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subdir, 0755)
	require.NoError(t, err)

	size, err := GetSize(tmpDir, true)
	require.NoError(t, err)
	assert.Equal(t, int64(16), size)
}

func TestGetSize_HiddenDirectory(t *testing.T) {
	parentDir, err := os.MkdirTemp("", "parent")
	require.NoError(t, err)
	defer os.RemoveAll(parentDir)

	hiddenDir := filepath.Join(parentDir, ".hiddendir")
	err = os.Mkdir(hiddenDir, 0755)
	require.NoError(t, err)

	hiddenFile := filepath.Join(hiddenDir, "file.txt")
	err = os.WriteFile(hiddenFile, []byte("content"), 0644)
	require.NoError(t, err)

	size, err := GetSize(parentDir, false)
	require.NoError(t, err)
	assert.Equal(t, int64(0), size)

	size, err = GetSize(parentDir, true)
	require.NoError(t, err)
	assert.Equal(t, int64(0), size)
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		human    bool
		expected string
	}{
		{"non-human small", 500, false, "500B"},
		{"non-human medium", 25165824, false, "25165824B"},

		{"human bytes", 123, true, "123B"},
		{"human KB", 1024, true, "1.0KB"},
		{"human KB fractional", 1536, true, "1.5KB"},
		{"human MB", 1048576, true, "1.0MB"},
		{"human MB fractional", 1572864, true, "1.5MB"},
		{"human GB", 1073741824, true, "1.0GB"},
		{"human TB", 1099511627776, true, "1.0TB"},
		{"human PB", 1125899906842624, true, "1.0PB"},
		{"human EB", 1152921504606846976, true, "1.0EB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, FormatSize(tt.size, tt.human))
		})
	}
}