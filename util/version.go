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

func VersionString() (string, error) {
	release, e := strconv.ParseBool(Release)
	if e != nil {
		return "", e
	}
	dirty, e := strconv.ParseBool(Dirty)
	if e != nil {
		return "", e
	}
	return versionString(Version, GitSha, release, dirty), nil
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
