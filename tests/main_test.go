package tests

import (
	"testing"
	"github.com/misterDillett/go-project-242/code"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSize(t *testing.T) {
	size, err := code.GetSize("testdata/file.txt", false, false)
	require.NoError(t, err)
	assert.Greater(t, size, int64(0))
}

func TestFormatSize(t *testing.T) {
	result := code.FormatSize(1024, true)
	assert.Equal(t, "1.0KB", result)
}