package fileio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileIORoundtrip(t *testing.T) {
	file, err := os.CreateTemp("", "ironclad-fileio-test")
	require.NoError(t, err, "failed to create temp fileio test file")
	defer os.Remove(file.Name())

	err = file.Close()
	require.NoError(t, err, "failed to close temp fileio test file")

	plaintext := []byte("foo bar baz 123 456 789")

	err = Save(file.Name(), "abc123", plaintext)
	require.NoError(t, err, "failed to save fileio test file")

	result, err := Load(file.Name(), "abc123")
	require.NoError(t, err, "failed to load fileio test file")

	require.Equal(t, plaintext, result)
}
