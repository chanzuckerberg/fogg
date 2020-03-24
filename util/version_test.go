package util

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	r := require.New(t)

	testCases := []struct {
		input   string
		version string
		sha     string
		dirty   bool
	}{
		{"0.1.0", "0.1.0", "", false},
		{"0.1.0-abcdef", "0.1.0", "abcdef", false},
		{"0.1.0-abcdef.dirty", "0.1.0", "abcdef", true},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			v, sha, dirty := ParseVersion(tc.input)
			semVersion, e := semver.Parse(tc.version)
			r.NoError(e)
			r.Equal(semVersion, v)
			r.Equal(tc.sha, sha)
			r.Equal(tc.dirty, dirty)
		})
	}
}
