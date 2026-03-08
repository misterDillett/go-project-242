// code_test.go
package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestStructure(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "test-*")
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "small.txt"), []byte("hello"), 0644) // 5 байт
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "medium.txt"), []byte(strings.Repeat("a", 2048)), 0644) // 2048 байт
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("secret"), 0644) // 6 байт
	require.NoError(t, err)

	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(subDir, "subfile.txt"), []byte("subcontent"), 0644) // 10 байт
	require.NoError(t, err)

	emptyDir := filepath.Join(tmpDir, "empty")
	err = os.Mkdir(emptyDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "привет.txt"), []byte("unicode"), 0644) // 7 байт
	require.NoError(t, err)

	return tmpDir
}

func TestGetPathSize_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "empty-*")
	require.NoError(t, err)
	defer func() {
		if rErr := os.Remove(tmpFile.Name()); rErr != nil {
			t.Logf("Warning: failed to remove temp file: %v", rErr)
		}
	}()

	result, err := GetPathSize(tmpFile.Name(), false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "0B", result)
}

func TestGetPathSize_SmallFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "small-*")
	require.NoError(t, err)
	defer func() {
		if rErr := os.Remove(tmpFile.Name()); rErr != nil {
			t.Logf("Warning: failed to remove temp file: %v", rErr)
		}
	}()

	content := []byte("hello")
	err = os.WriteFile(tmpFile.Name(), content, 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpFile.Name(), false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "5B", result)

	result, err = GetPathSize(tmpFile.Name(), false, true, false)
	require.NoError(t, err)
	assert.Equal(t, "5B", result)
}

func TestGetPathSize_EmptyDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "empty-*")
	require.NoError(t, err)
	defer func() {
		if rmErr := os.RemoveAll(tmpDir); rmErr != nil {
			t.Logf("Warning: failed to remove temp dir: %v", rmErr)
		}
	}()

	result, err := GetPathSize(tmpDir, false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "0B", result)
}

func TestGetPathSize_DirectoryNonRecursive(t *testing.T) {
	tmpDir := createTestStructure(t)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	result, err := GetPathSize(tmpDir, false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "2060B", result)

	result, err = GetPathSize(tmpDir, false, false, true)
	require.NoError(t, err)
	assert.Equal(t, "2066B", result)
}

func TestGetPathSize_DirectoryRecursive(t *testing.T) {
	tmpDir := createTestStructure(t)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	result, err := GetPathSize(tmpDir, true, false, false)
	require.NoError(t, err)
	assert.Equal(t, "2070B", result)

	result, err = GetPathSize(tmpDir, true, false, true)
	require.NoError(t, err)
	assert.Equal(t, "2076B", result)
}

func TestGetPathSize_HumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{"Bytes", 500, "500B"},
		{"KB", 1024, "1.0KB"},
		{"KB fractional", 1536, "1.5KB"},
		{"MB", 1048576, "1.0MB"},
		{"MB fractional", 1572864, "1.5MB"},
		{"GB", 1073741824, "1.0GB"},
		// {"GB fractional", 1610612736, "1.5GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp(t.TempDir(), "test-*")
			require.NoError(t, err)

			err = os.WriteFile(tmpFile.Name(), make([]byte, tt.size), 0644)
			require.NoError(t, err)

			result, err := GetPathSize(tmpFile.Name(), false, true, false)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPathSize_UnicodePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "тест-*")
	require.NoError(t, err)
	defer func() {
		if rmErr := os.RemoveAll(tmpDir); rmErr != nil {
			t.Logf("Warning: failed to remove temp dir: %v", rmErr)
		}
	}()

	filePath := filepath.Join(tmpDir, "привет-мир.txt")
	err = os.WriteFile(filePath, []byte("unicode"), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(filePath, false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "7B", result)

	result, err = GetPathSize(tmpDir, false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "7B", result)
}

func TestGetPathSize_Symlinks(t *testing.T) {
	tmpDir := t.TempDir()

	targetFile := filepath.Join(tmpDir, "target.txt")
	err := os.WriteFile(targetFile, []byte("target"), 0644) // 6 байт
	require.NoError(t, err)

	symlinkPath := filepath.Join(tmpDir, "link.txt")
	err = os.Symlink(targetFile, symlinkPath)
	if err != nil {
		t.Skip("Symlinks not supported on this OS")
	}

	t.Run("symlink file itself", func(t *testing.T) {
		result, err := GetPathSize(symlinkPath, false, false, false)
		require.NoError(t, err)

		assert.NotEmpty(t, result)
		t.Logf("Symlink size: %s", result)
	})

	t.Run("directory with symlink non-recursive", func(t *testing.T) {
		result, err := GetPathSize(tmpDir, false, false, false)
		require.NoError(t, err)

		t.Logf("Directory size with symlink: %s", result)
		assert.NotEmpty(t, result)
	})
}

func TestGetPathSize_NonExistentPath(t *testing.T) {
	_, err := GetPathSize("/path/does/not/exist", false, false, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "/path/does/not/exist")
}
