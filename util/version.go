package util

import (
	"fmt"
	"strconv"
)

var (
	Version string
	GitSha  string
	Release string
	Dirty   string
)

func VersionString() string {
	release, _ := strconv.ParseBool(Release)
	dirty, _ := strconv.ParseBool(Dirty)
	return versionString(Version, GitSha, release, dirty)
}

func versionString(version, sha string, release, dirty bool) string {
	if release {
		return version
	}
	if !dirty {
		return fmt.Sprintf("%s-%s", version, sha)
	}
	return fmt.Sprintf("%s-%s-dirty", version, sha)
}
