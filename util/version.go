package util

import (
	"fmt"
	"strconv"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

var (
	Version = "undefined"
	GitSha  = "undefined"
	Release = "false"
	Dirty   = "true"
)

func VersionString() (string, error) {
	release, e := strconv.ParseBool(Release)
	if e != nil {
		return "", errors.Wrapf(e, "unable to parse version release field %s", Release)
	}
	dirty, e := strconv.ParseBool(Dirty)
	if e != nil {
		return "", errors.Wrapf(e, "unable to parse version dirty field %s", Dirty)
	}
	return versionString(Version, GitSha, release, dirty), nil
}

func VersionCacheKey() string {
	versionString, e := VersionString()
	if e != nil {
		return ""
	}
	v, e := semver.Parse(versionString)
	if e != nil {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
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
