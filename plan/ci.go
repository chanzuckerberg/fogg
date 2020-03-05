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
	projects    []CIProject

	TestBuckets [][]CIProject
	AWSProfiles ciAwsProfiles
	Buildevents bool
}

type CircleCIConfig struct {
	CIConfig
	SSHFingerprints []string
}

type GitHubActionsCIConfig struct {
	CIConfig
}

type TravisCIConfig struct {
	CIConfig
}

func (c *CIConfig) populateBuckets(numBuckets int) *CIConfig {
	sort.SliceStable(c.projects, func(i, j int) bool {
		return c.projects[i].Name < c.projects[j].Name
	})

	c.TestBuckets = make([][]CIProject, numBuckets)
	for i, proj := range c.projects {
		bucket := i % numBuckets
		c.TestBuckets[bucket] = append(c.TestBuckets[bucket], proj)
	}
	return c
}

func (c *CIConfig) addProjects(projects ...CIProject) *CIConfig {
	if c == nil {
		c = &CIConfig{}
	}

	c.projects = append(c.projects, projects...)
	return c
}

func (c *CIConfig) merge(other *CIConfig) *CIConfig {
	if c == nil {
		c = &CIConfig{}
	}
	if other == nil {
		return c
	}

	c.Enabled = c.Enabled || other.Enabled
	c.Buildevents = c.Buildevents || other.Buildevents
	c.AWSProfiles = c.AWSProfiles.merge(other.AWSProfiles)
	c.projects = append(c.projects, other.projects...)

	return c
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

func (p *Plan) buildTravisCIConfig(c *v2.Config, foggVersion string) TravisCIConfig {
	ciConfig := &CIConfig{
		FoggVersion: foggVersion,
	}

	globalConfig := p.Global.TravisCI.generateCIConfig(
		p.Global.Backend,
		p.Global.Providers.AWS,
		"global",
		"terraform/global")
	ciConfig = ciConfig.merge(globalConfig)

	for name, acct := range p.Accounts {
		accountConfig := acct.TravisCI.generateCIConfig(
			acct.Backend,
			acct.Providers.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)
		ciConfig = ciConfig.merge(accountConfig)
	}

	for envName, env := range p.Envs {
		for cName, c := range env.Components {
			envConfig := c.TravisCI.generateCIConfig(
				c.Backend,
				c.Providers.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)
			ciConfig = ciConfig.merge(envConfig)
		}
	}

	for moduleName := range p.Modules {
		proj := CIProject{
			Name:    fmt.Sprintf("modules/%s", moduleName),
			Dir:     fmt.Sprintf("terraform/modules/%s", moduleName),
			Command: "check",
		}
		ciConfig = ciConfig.addProjects(proj)
	}

	numBuckets := 1
	if c.Defaults.Tools != nil &&
		c.Defaults.Tools.TravisCI != nil {

		if c.Defaults.Tools.TravisCI.TestBuckets != nil &&
			*c.Defaults.Tools.TravisCI.TestBuckets > 0 {
			numBuckets = *c.Defaults.Tools.TravisCI.TestBuckets
		}

		// If aws is disabled, reset the providers
		aws, ok := c.Defaults.Tools.TravisCI.Providers["aws"]
		if ok && aws.Disabled {
			ciConfig.AWSProfiles = ciAwsProfiles{}
		}
	}

	ciConfig = ciConfig.populateBuckets(numBuckets)
	return TravisCIConfig{
		CIConfig: *ciConfig,
	}
}

func (p *Plan) buildCircleCIConfig(c *v2.Config, foggVersion string) CircleCIConfig {
	ciConfig := &CIConfig{
		FoggVersion: foggVersion,
	}

	globalConfig := p.Global.CircleCI.generateCIConfig(
		p.Global.Backend,
		p.Global.Providers.AWS,
		"global",
		"terraform/global")
	ciConfig = ciConfig.merge(globalConfig)

	for name, acct := range p.Accounts {
		accountConfig := acct.CircleCI.generateCIConfig(
			acct.Backend,
			acct.Providers.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)
		ciConfig = ciConfig.merge(accountConfig)
	}

	for envName, env := range p.Envs {
		for cName, c := range env.Components {
			envConfig := c.CircleCI.generateCIConfig(
				c.Backend,
				c.Providers.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)
			ciConfig = ciConfig.merge(envConfig)
		}
	}

	for moduleName := range p.Modules {
		proj := CIProject{
			Name:    fmt.Sprintf("modules/%s", moduleName),
			Dir:     fmt.Sprintf("terraform/modules/%s", moduleName),
			Command: "check",
		}
		ciConfig = ciConfig.addProjects(proj)
	}

	numBuckets := 1
	sshFingerprints := []string{}

	if c.Defaults.Tools != nil && c.Defaults.Tools.CircleCI != nil {

		if c.Defaults.Tools.CircleCI.TestBuckets != nil &&
			*c.Defaults.Tools.CircleCI.TestBuckets > 0 {
			numBuckets = *c.Defaults.Tools.CircleCI.TestBuckets
		}

		sshFingerprints = append(sshFingerprints, c.Defaults.Tools.CircleCI.SSHKeyFingerprints...)

		// If aws is disabled, reset the providers
		aws, ok := c.Defaults.Tools.CircleCI.Providers["aws"]
		if ok && aws.Disabled {
			ciConfig.AWSProfiles = ciAwsProfiles{}
		}
	}

	ciConfig = ciConfig.populateBuckets(numBuckets)
	return CircleCIConfig{
		CIConfig:        *ciConfig,
		SSHFingerprints: sshFingerprints,
	}
}

func (p *Plan) buildGitHubActionsConfig(c *v2.Config, foggVersion string) GitHubActionsCIConfig {
	ciConfig := &CIConfig{
		FoggVersion: foggVersion,
	}

	globalConfig := p.Global.GitHubActionsCI.generateCIConfig(
		p.Global.Backend,
		p.Global.Providers.AWS,
		"global",
		"terraform/global")
	ciConfig = ciConfig.merge(globalConfig)

	for name, acct := range p.Accounts {
		accountConfig := acct.GitHubActionsCI.generateCIConfig(
			acct.Backend,
			acct.Providers.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)
		ciConfig = ciConfig.merge(accountConfig)
	}

	for envName, env := range p.Envs {
		for cName, c := range env.Components {
			envConfig := c.GitHubActionsCI.generateCIConfig(
				c.Backend,
				c.Providers.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)
			ciConfig = ciConfig.merge(envConfig)
		}
	}

	for moduleName := range p.Modules {
		proj := CIProject{
			Name:    fmt.Sprintf("modules/%s", moduleName),
			Dir:     fmt.Sprintf("terraform/modules/%s", moduleName),
			Command: "check",
		}
		ciConfig = ciConfig.addProjects(proj)
	}

	numBuckets := 1
	if c.Defaults.Tools != nil && c.Defaults.Tools.GitHubActionsCI != nil {

		if c.Defaults.Tools.GitHubActionsCI.TestBuckets != nil &&
			*c.Defaults.Tools.GitHubActionsCI.TestBuckets > 0 {
			numBuckets = *c.Defaults.Tools.GitHubActionsCI.TestBuckets
		}

		// If aws is disabled, reset the providers
		aws, ok := c.Defaults.Tools.GitHubActionsCI.Providers["aws"]
		if ok && aws.Disabled {
			ciConfig.AWSProfiles = ciAwsProfiles{}
		}
	}

	ciConfig = ciConfig.populateBuckets(numBuckets)
	return GitHubActionsCIConfig{
		CIConfig: *ciConfig,
	}
}
