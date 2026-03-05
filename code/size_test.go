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

	size, err := GetSize(tmpFile.Name(), false, false)
	require.NoError(t, err)
	assert.Equal(t, int64(len(content)), size)
}

func TestGetSize_Directory_NonRecursive(t *testing.T) {

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

	subfile := filepath.Join(subdir, "subfile.txt")
	err = os.WriteFile(subfile, []byte("subcontent"), 0644)
	require.NoError(t, err)

	size, err := GetSize(tmpDir, false, false)
	require.NoError(t, err)
	assert.Equal(t, int64(10), size)

	size, err = GetSize(tmpDir, false, true)
	require.NoError(t, err)
	assert.Equal(t, int64(16), size)
}

func TestGetSize_Directory_Recursive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.txt")
	err = os.WriteFile(file1, []byte("hello"), 0644)
	require.NoError(t, err)

	hiddenFile := filepath.Join(tmpDir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("secret"), 0644)
	require.NoError(t, err)

	subdir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subdir, 0755)
	require.NoError(t, err)

	subfile := filepath.Join(subdir, "subfile.txt")
	err = os.WriteFile(subfile, []byte("subcontent"), 0644)
	require.NoError(t, err)

	nestedDir := filepath.Join(subdir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	require.NoError(t, err)

	nestedFile := filepath.Join(nestedDir, "nested.txt")
	err = os.WriteFile(nestedFile, []byte("nested"), 0644)
	require.NoError(t, err)

	size, err := GetSize(tmpDir, true, false)
	require.NoError(t, err)
	expected := int64(5 + 10 + 6)
	assert.Equal(t, expected, size)

	size, err = GetSize(tmpDir, true, true)
	require.NoError(t, err)
	expected += 6
	assert.Equal(t, expected, size)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, FormatSize(tt.size, tt.human))
		})
	}
}