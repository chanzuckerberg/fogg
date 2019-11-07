package plan

import (
	"fmt"
	"sort"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
)

type CIProject struct {
	Name    string
	Dir     string
	Command string
}
type CIConfig struct {
	Enabled     bool
	FoggVersion string
	TestBuckets [][]CIProject
	AWSProfiles map[string]AWSRole
	Buildevents bool
}

type ciAwsProfiles map[string]AWSRole

func (p ciAwsProfiles) merge(other ciAwsProfiles) ciAwsProfiles {
	if p == nil {
		p = ciAwsProfiles{}
	}

	for profile, role := range other {
		p[profile] = role
	}

	return p
}

// TODO(el): mostly a duplicate of buildAtlantis(). refactor later
func (p *Plan) buildTravisCIConfig(c *v2.Config, foggVersion string) CIConfig {
	enabled := false
	buildeventsEnabled := false
	projects := []CIProject{}
	awsProfiles := ciAwsProfiles{}

	globalProjects, globalProfiles, globalBuildeventsEnabled := p.Global.TravisCI.generateCIConfig(
		p.Global.Backend,
		p.Global.Providers.AWS,
		"global",
		"terraform/global")
	buildeventsEnabled = buildeventsEnabled || globalBuildeventsEnabled
	projects = append(projects, globalProjects...)
	awsProfiles.merge(globalProfiles)

	for name, acct := range p.Accounts {
		accountProjects, accountProfiles, accoutnBuildeventsEnabled := acct.TravisCI.generateCIConfig(
			acct.Backend,
			acct.Providers.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)
		buildeventsEnabled = buildeventsEnabled || accoutnBuildeventsEnabled
		projects = append(projects, accountProjects...)
		awsProfiles.merge(accountProfiles)
		enabled = enabled || len(projects) > 0
	}

	for envName, env := range p.Envs {
		for cName, c := range env.Components {
			envProjects, envProfiles, envBuildeventsEnabled := c.TravisCI.generateCIConfig(
				c.Backend,
				c.Providers.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)
			buildeventsEnabled = buildeventsEnabled || envBuildeventsEnabled
			projects = append(projects, envProjects...)
			awsProfiles.merge(envProfiles)
			enabled = enabled || len(projects) > 0
		}
	}

	for moduleName := range p.Modules {
		proj := CIProject{
			Name:    fmt.Sprintf("modules/%s", moduleName),
			Dir:     fmt.Sprintf("terraform/modules/%s", moduleName),
			Command: "check",
		}
		projects = append(projects, proj)
	}

	buckets := 1
	if c.Defaults.Tools != nil &&
		c.Defaults.Tools.TravisCI != nil &&
		c.Defaults.Tools.TravisCI.TestBuckets != nil &&
		*c.Defaults.Tools.TravisCI.TestBuckets > 0 {

		buckets = *c.Defaults.Tools.TravisCI.TestBuckets
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	testBuckets := make([][]CIProject, buckets)
	for i, proj := range projects {
		bucket := i % buckets
		testBuckets[bucket] = append(testBuckets[bucket], proj)
	}

	tr := CIConfig{
		Enabled:     enabled,
		Buildevents: buildeventsEnabled,
		FoggVersion: foggVersion,
		TestBuckets: testBuckets,
		AWSProfiles: awsProfiles,
	}
	return tr
}
