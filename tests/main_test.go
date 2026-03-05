package tests

import (
    "testing"
    "hexlet-boilerplates/gopackage/code"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetSize(t *testing.T) {
    size, err := code.GetSize("testdata/file.txt", false, false)
    require.NoError(t, err)
    assert.Greater(t, size, int64(0))
}