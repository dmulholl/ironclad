package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	err := Set("foo", "abc")
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
