package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/chanzuckerberg/fogg/errs"
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
		return "", errs.WrapInternal(e, fmt.Sprintf("unable to parse version release field %s", Release))
	}
	dirty, e := strconv.ParseBool(Dirty)
	if e != nil {
		return "", errs.WrapInternal(e, fmt.Sprintf("unable to parse version dirty field %s", Dirty))
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

func ParseVersion(version string) (semver.Version, string, bool) {
	var dirty bool
	var sha string
	v := version
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

func versionString(version, sha string, release, dirty bool) string {
	if release {
		return version
	}
	if !dirty {
		return fmt.Sprintf("%s-%s", version, sha)
	}
	return fmt.Sprintf("%s-%s.dirty", version, sha)
}
