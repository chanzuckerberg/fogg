package plan

import (
	"encoding/json"
	"path"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
)

//AWSProfile represents Travis CI's AWS profile
type AWSProfile struct {
	Name string
	ID   json.Number
	Role string
}

//TravisCI represents the travis ci configuration
type TravisCI struct {
	AWSProfiles []AWSProfile
	Enabled     bool
	FoggVersion string
	TestBuckets [][]string
}

func (p *Plan) buildTravisCI(c *v2.Config, version string) TravisCI {
	if p.Accounts == nil {
		panic("buildTravisCI must be run after buildAccounts")
	}

	if c.Defaults.Tools == nil || c.Defaults.Tools.TravisCI == nil {
		return TravisCI{}
	}

	tr := TravisCI{
		Enabled: c.Defaults.Tools.TravisCI.Enabled,
	}
	var profiles []AWSProfile

	tr.FoggVersion = version

	for _, name := range util.SortedMapKeys(p.Accounts) {
		profiles = append(profiles, AWSProfile{
			Name: name,
			// TODO since accountID is required here, that means we need
			// to make it non-optional, either in defaults or post-plan.
			ID:   p.Accounts[name].Providers.AWS.AccountID,
			Role: c.Defaults.Tools.TravisCI.AWSIAMRoleName,
		})
	}
	tr.AWSProfiles = profiles

	var buckets int
	if c.Defaults.Tools.TravisCI.TestBuckets > 0 {
		buckets = c.Defaults.Tools.TravisCI.TestBuckets
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
