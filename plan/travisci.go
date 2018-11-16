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
	if p.Accounts == nil {
		panic("buildTravisCI must be run after buildAccounts")
	}

	tr := TravisCI{
		Enabled:          c.TravisCI.Enabled,
		AWSIDAccountName: c.TravisCI.IDAccountName,
	}
	var profiles []AWSProfile

	for name, a := range p.Accounts {
		profiles = append(profiles, AWSProfile{
			Name: name,
			// TODO since accountID is required here, that means we need
			// to make it non-optional, either in defaults or post-plan.
			ID:            *a.AccountID,
			Role:          c.TravisCI.AWSIAMRoleName,
			IDAccountName: c.TravisCI.IDAccountName,
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
	jobs := make(chan string)
	done := make(chan bool)

	go func() {
		var componentCount int

		for {
			j, more := <-jobs
			if more {
				shard := componentCount % shards
				testShards[shard] = append(testShards[shard], j)
				componentCount++
			} else {
				done <- true
				return
			}
		}
	}()

	jobs <- path.Join("terraform", "global")

	for _, name := range util.SortedMapKeys(c.Accounts) {
		jobs <- path.Join("terraform", "accounts", name)
	}

	for _, envName := range util.SortedMapKeys(c.Envs) {
		for _, name := range util.SortedMapKeys(c.Envs[envName].Components) {
			jobs <- path.Join("terraform", "envs", envName, name)
		}
	}

	for _, moduleName := range util.SortedMapKeys(c.Modules) {
		jobs <- path.Join("terraform", "modules", moduleName)
	}

	close(jobs)
	<-done
	tr.TestShards = testShards
	return tr

}
