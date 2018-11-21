package github

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	version "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Versions represents github release versions
type Versions struct {
	versions []*version.Version
}

// Latest resturns the latest version if one exists, else zero value
func (v *Versions) Latest() string {
	current := v.latest()
	if current == nil {
		return ""
	}
	return current.String()
}

// latest gets the latest version
func (v *Versions) latest() *version.Version {
	if len(v.versions) < 1 {
		return nil
	}
	return v.versions[0]
}

// Outdated returns true if there is a newer version
func (v *Versions) Outdated(ver string) (bool, error) {
	latest := v.latest()
	if latest == nil {
		return false, nil
	}
	testVersion, err := version.NewVersion(ver)
	if err != nil {
		return false, errors.Wrapf(err, "Could not parse %s", ver)
	}
	return latest.GreaterThan(testVersion), nil
}

// CheckLatestVersion checks to see if we're on the latest version
func (c *Client) CheckLatestVersion(
	ctx context.Context,
	repoOwner, repoName, currentVersion string) (*Versions, error) {

	// page through all releases
	allReleases := []*github.RepositoryRelease{}
	pageOptions := &github.ListOptions{}
	for {
		releases, resp, err := c.client.Repositories.ListReleases(ctx, repoOwner, repoName, pageOptions)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not fetch releases")
		}
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("Unknown status code %d", resp.StatusCode)
		}
		allReleases = append(allReleases, releases...)
		if resp.NextPage == 0 {
			break
		}
		pageOptions.Page = resp.NextPage
	}

	// get all valid rawVersions
	versions := []*version.Version{}
	for _, release := range allReleases {
		if release != nil && release.TagName != nil {
			// trim a leading v
			trimmedRelease := strings.TrimPrefix(release.GetTagName(), "v")

			v, err := version.NewVersion(trimmedRelease)
			if err != nil {
				logrus.WithError(err).Warnf("Could not semver parse tag %s", trimmedRelease)
				continue
			}
			versions = append(versions, v)
		}
	}

	// sort the versions
	sort.Sort(version.Collection(versions))
	return &Versions{versions: versions}, nil
}
