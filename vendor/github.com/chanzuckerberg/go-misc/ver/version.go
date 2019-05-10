package ver

import (
	"fmt"
	"strconv"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

// VersionString returns the version string
func VersionString(version, gitsha, releaseStr, dirtyStr string) (string, error) {
	release, e := strconv.ParseBool(releaseStr)
	if e != nil {
		return "", errors.Wrapf(e, "unable to parse version release field %s", releaseStr)
	}
	dirty, e := strconv.ParseBool(dirtyStr)
	if e != nil {
		return "", errors.Wrapf(e, "unable to parse version dirty field %s", dirtyStr)
	}
	return versionString(version, gitsha, release, dirty), nil
}

// VersionCacheKey returns a key to version the cache
func VersionCacheKey(version, gitsha, release, dirty string) string {
	versionString, e := VersionString(version, gitsha, release, dirty)
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
		return fmt.Sprintf("%s+%s", version, sha)
	}
	return fmt.Sprintf("%s+%s+dirty", version, sha)
}
