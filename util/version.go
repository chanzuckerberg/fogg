package util

import (
	"strings"

	"github.com/blang/semver"
	"github.com/chanzuckerberg/go-misc/ver"
)

var (
	Version = "undefined"
	GitSha  = "undefined"
	Release = "false"
	Dirty   = "true"
)

func VersionString() (string, error) {
	return ver.VersionString(Version, GitSha, Release, Dirty)
}

func VersionCacheKey() string {
	return ver.VersionCacheKey(Version, GitSha, Release, Dirty)
}

func ParseVersion(version string) (semver.Version, string, bool) {
	var dirty bool
	var sha string
	v := strings.TrimSpace(version)
	if strings.HasSuffix(v, ".dirty") {
		dirty = true
		v = strings.TrimSuffix(v, ".dirty")
	}
	if strings.Contains(v, "-") {
		tmp := strings.Split(v, "-")
		v = tmp[0]
		sha = tmp[1]
	}

	semVersion, _ := semver.Parse(v)
	return semVersion, sha, dirty
}
