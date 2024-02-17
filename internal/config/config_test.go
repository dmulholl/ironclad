package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	file, err := os.CreateTemp("", "ironclad-config-test")
	require.NoError(t, err, "failed to create temp config file")
	defer os.Remove(file.Name())

	err = file.Close()
	require.NoError(t, err, "failed to close temp config file")

	err = os.Remove(file.Name())
	require.NoError(t, err, "failed to delete temp config file")

	ConfigFile = file.Name()

	// Try setting to a config file that doesn't exist.
	err = Set("foo", "abc")
	require.NoError(t, err, "failed to set 'foo=abc'")

	value, found, err := Get("foo")
	require.NoError(t, err, "failed to get 'foo'")
	require.True(t, found, "failed to find 'foo'")
	require.Equal(t, value, "abc")

	err = Delete("foo")
	require.NoError(t, err, "failed to delete 'foo'")

	_, found, err = Get("foo")
	require.NoError(t, err, "failed to get 'foo' after deleting")
	require.False(t, found, "found 'foo' after deleting")

	// Try setting to a config file that does exist.
	err = Set("bar", "def")
	require.NoError(t, err, "failed to set 'bar=def'")

	value, found, err = Get("bar")
	require.NoError(t, err, "failed to get 'bar'")
	require.True(t, found, "failed to find 'bar'")
	require.Equal(t, value, "def")

	err = Delete("bar")
	require.NoError(t, err, "failed to delete 'bar'")

	_, found, err = Get("bar")
	require.NoError(t, err, "failed to get 'bar' after deleting")
	require.False(t, found, "found 'bar' after deleting")
}
