package util

import (
	"context"
	"fmt"

	"github.com/chanzuckerberg/go-misc/github"
	"github.com/sirupsen/logrus"
)

// CheckLatestVersion checks to see if we're on the latest version
func CheckLatestVersion(ctx context.Context, owner string, repo string) error {
	err := checkLatestVersion(ctx, owner, repo)
	if err != nil {
		logrus.WithError(err).WithField("repo", fmt.Sprintf("%s/%s", owner, repo)).Warn("Could not fetch latest release info")
		return err
	}
	return nil
}

func checkLatestVersion(ctx context.Context, owner string, repo string) error {
	currentVersion, err := VersionString()
	if err != nil {
		return err
	}
	githubClient := github.NewClient(nil)
	versions, err := githubClient.CheckLatestVersion(ctx, owner, repo, currentVersion)
	if err != nil {
		return err
	}

	outdated, err := versions.Outdated(currentVersion)
	if err != nil {
		return err
	}

	if outdated {
		logrus.
			WithField("current_version", currentVersion).
			WithField("latest_version", versions.Latest()).
			Warnf("Please upgrade %s/%s", owner, repo)
	}
	return nil
}
