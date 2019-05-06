package plan

import (
	"path"

	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
)

type TravisCI struct {
	AWSProfiles []AWSProfile
	Docker      bool
	Enabled     bool
	FoggVersion string
	TestBuckets [][]string
}

func (p *Plan) buildTravisCI(c *v2.Config, version string) TravisCI {
	if p.Accounts == nil {
		panic("buildTravisCI must be run after buildAccounts")
	}

	tr := TravisCI{
		Enabled: c.Tools.TravisCI.Enabled,
	}
	var profiles []AWSProfile

	tr.Docker = c.Docker
	tr.FoggVersion = version

	for _, name := range util.SortedMapKeys(p.Accounts) {
		profiles = append(profiles, AWSProfile{
			Name: name,
			// TODO since accountID is required here, that means we need
			// to make it non-optional, either in defaults or post-plan.
			ID:   p.Accounts[name].AccountID,
			Role: c.Tools.TravisCI.AWSIAMRoleName,
		})
	}
	tr.AWSProfiles = profiles

	var buckets int
	if c.Tools.TravisCI.TestBuckets > 0 {
		buckets = c.Tools.TravisCI.TestBuckets
	} else {
		buckets = 1
	}

	var testPaths []string
	testPaths = append(testPaths, path.Join("terraform", "global"))

	for _, name := range util.SortedMapKeys(c.Accounts) {
		testPaths = append(testPaths, path.Join("terraform", "accounts", name))
	}

	for _, envName := range util.SortedMapKeys(c.Envs) {
		for _, name := range util.SortedMapKeys(c.Envs[envName].Components) {
			testPaths = append(testPaths, path.Join("terraform", "envs", envName, name))
		}
	}

	for _, moduleName := range util.SortedMapKeys(c.Modules) {
		testPaths = append(testPaths, path.Join("terraform", "modules", moduleName))
	}

	testBuckets := make([][]string, buckets)
	for i, path := range testPaths {
		bucket := i % buckets
		testBuckets[bucket] = append(testBuckets[bucket], path)
	}

	tr.TestBuckets = testBuckets
	return tr
}
