package util

import (
	"fmt"
	"regexp"
)

const AWS_CODEARTIFACT_REGISTRY_REGEX = `\.codeartifact.*\.amazonaws\.com`
const AWS_CODEARTIFACT_CAPTURE_REGEX = `([a-z0-9-]+)-(.+)\.d\.codeartifact\.(.+)\.amazonaws\.com\/.*/([^\/.]+)\/`

// Evaluates if the `registryUrl` is a AWS CodeArtifact registry.
// Returns true if the URL is a AWS CodeArtifact registry, false otherwise.
func IsCodeArtifactURL(registryUrl string) bool {
	// ref: https://github.com/projen/projen/blob/v0.91.4/src/release/publisher.ts#L1252
	regex := regexp.MustCompile(AWS_CODEARTIFACT_REGISTRY_REGEX)
	return regex.MatchString(registryUrl)
}

type CodeArtifactRepository struct {
	Domain    string
	AccountId string
	Region    string
	Name      string
}

// Gets AWS details from the Code Artifact `registryUrl`.
// throws exception if not matching expected pattern
func ParseRegistryUrl(registryUrl string) (*CodeArtifactRepository, error) {
	// https://github.com/projen/projen/blob/v0.91.4/src/javascript/util.ts#L48
	regex := regexp.MustCompile(AWS_CODEARTIFACT_CAPTURE_REGEX)
	matches := regex.FindStringSubmatch(registryUrl)
	if len(matches) == 0 {
		return nil, fmt.Errorf("registry URL is not a valid CodeArtifact Repository URL, got: %s", registryUrl)
	}
	return &CodeArtifactRepository{
		Domain:    matches[1],
		AccountId: matches[2],
		Region:    matches[3],
		Name:      matches[4],
	}, nil
}
