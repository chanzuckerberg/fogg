package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	m := make(map[string]string)
	m["foo"] = "bar"
	r := dict(m)
	require.NotNil(t, r)
	require.IsType(t, map[string]interface{}{}, r)
	require.Equal(t, "bar", r["foo"])
}
