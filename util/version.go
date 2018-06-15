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

	if release {
		return Version
	}
	if !dirty {
		return fmt.Sprintf("%s-%s", Version, GitSha)
	}
	return fmt.Sprintf("%s-%s-dirty", Version, GitSha)
}
