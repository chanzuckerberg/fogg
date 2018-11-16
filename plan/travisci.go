package plan

import (
	"path"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
)

type TravisCI struct {
	Enabled          bool
	AWSIDAccountName string
	AWSProfiles      []AWSProfile
	TestShards       [][]string
}

func (p *Plan) buildTravisCI(c *config.Config) TravisCI {
	tr := TravisCI{
		Enabled:          c.TravisCI.Enabled,
		AWSIDAccountName: c.TravisCI.HubAccountName,
	}
	var profiles []AWSProfile
	// TODO we should actually take the resolved values for these, not
	//  raw config
	for name, a := range c.Accounts {
		profiles = append(profiles, AWSProfile{
			Name:           name,
			ID:             *a.AccountID,
			Role:           c.TravisCI.AWSIAMRoleName,
			HubAccountName: c.TravisCI.HubAccountName,
		})
	}
	tr.AWSProfiles = profiles

	var shards int
	if c.TravisCI.TestShards > 0 {
		shards = c.TravisCI.TestShards
	} else {
		shards = 1
	}

	testShards := make([][]string, shards)

	var componentCount int

	shard := componentCount % shards
	testShards[shard] = append(testShards[shard], path.Join("terraform", "global"))
	componentCount++

	for _, name := range util.SortedMapKeys(c.Accounts) {
		shard := componentCount % shards
		testShards[shard] = append(testShards[shard], path.Join("terraform", "accounts", name))
		componentCount++
	}

	for _, envName := range util.SortedMapKeys(c.Envs) {
		for _, name := range util.SortedMapKeys(c.Envs[envName].Components) {
			shard := componentCount % shards
			testShards[shard] = append(testShards[shard], path.Join("terraform", "envs", envName, name))
			componentCount++

		}
	}

	tr.TestShards = testShards
	return tr

}
