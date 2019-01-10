package util

import (
	"testing"

	masterminds_semver "github.com/Masterminds/semver"
	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
)

func TestVersionString(t *testing.T) {
	s := versionString("0.1.0", "abcdef", true, false)
	assert.Equal(t, "0.1.0", s)
	assert.Equal(t, "", masterminds_semver.MustParse(s).Prerelease())

	s = versionString("0.1.0", "abcdef", false, false)
	assert.Equal(t, "0.1.0-abcdef", s)
	assert.Equal(t, "abcdef", masterminds_semver.MustParse(s).Prerelease())

	s = versionString("0.1.0", "abcdef", false, true)
	assert.Equal(t, "0.1.0-abcdef.dirty", s)
	assert.Equal(t, "abcdef.dirty", masterminds_semver.MustParse(s).Prerelease())

}

func TestParse(t *testing.T) {
	a := assert.New(t)

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
			a.NoError(e)
			a.Equal(semVersion, v)
			a.Equal(tc.sha, sha)
			a.Equal(tc.dirty, dirty)
		})
	}

}
