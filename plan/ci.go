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
	Env         map[string]string
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
	SSHKeySecrets []string
	RunsOn        []string
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

func (c *CIConfig) merge(other *CIConfig, awsProvider v2.CIProviderConfig) *CIConfig {
	if c == nil {
		c = &CIConfig{}
	}
	if other == nil {
		return c
	}

	c.Enabled = c.Enabled || other.Enabled
	c.Buildevents = c.Buildevents || other.Buildevents
	c.AWSProfiles = c.AWSProfiles.merge(other.AWSProfiles, awsProvider)
	c.projects = append(c.projects, other.projects...)

	return c
}

type ciAwsProfiles map[string]AWSRole

func (p ciAwsProfiles) merge(other ciAwsProfiles, awsProvider v2.CIProviderConfig) ciAwsProfiles {
	if p == nil {
		p = ciAwsProfiles{}
	}

	if awsProvider.Disabled {
		return p
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

	var awsProvider v2.CIProviderConfig

	if c.Defaults.Tools != nil && c.Defaults.Tools.TravisCI != nil {
		awsProvider = c.Defaults.Tools.TravisCI.Providers["aws"]
	}

	globalConfig := p.Global.TravisCI.generateCIConfig(
		p.Global.Backend,
		p.Global.ProviderConfiguration.AWS,
		"global",
		"terraform/global")

	if c.Global.Tools != nil && c.Global.Tools.TravisCI != nil {
		ciConfig = ciConfig.merge(globalConfig, c.Global.Tools.TravisCI.Providers["aws"])
	} else {
		ciConfig = ciConfig.merge(globalConfig, awsProvider)
	}

	for name, acct := range p.Accounts {
		accountConfig := acct.TravisCI.generateCIConfig(
			acct.Backend,
			acct.ProviderConfiguration.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)

		if c.Accounts[name].Tools != nil && c.Accounts[name].Tools.TravisCI != nil {
			ciConfig = ciConfig.merge(accountConfig, c.Accounts[name].Tools.TravisCI.Providers["aws"])
		} else {
			ciConfig = ciConfig.merge(accountConfig, awsProvider)
		}
	}

	for envName, env := range p.Envs {
		for cName, d := range env.Components {
			envConfig := d.TravisCI.generateCIConfig(
				d.Backend,
				d.ProviderConfiguration.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)

			if c.Envs[envName].Tools != nil && c.Envs[envName].Tools.TravisCI != nil {
				ciConfig = ciConfig.merge(envConfig, c.Envs[envName].Tools.TravisCI.Providers["aws"])
			} else {
				ciConfig = ciConfig.merge(envConfig, awsProvider)
			}
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

	var awsProvider v2.CIProviderConfig

	if c.Defaults.Tools != nil && c.Defaults.Tools.CircleCI != nil {
		awsProvider = c.Defaults.Tools.CircleCI.Providers["aws"]
	}

	globalConfig := p.Global.CircleCI.generateCIConfig(
		p.Global.Backend,
		p.Global.ProviderConfiguration.AWS,
		"global",
		"terraform/global")

	if c.Global.Tools != nil && c.Global.Tools.CircleCI != nil {
		ciConfig = ciConfig.merge(globalConfig, c.Global.Tools.CircleCI.Providers["aws"])
	} else {
		ciConfig = ciConfig.merge(globalConfig, awsProvider)
	}

	for name, acct := range p.Accounts {
		accountConfig := acct.CircleCI.generateCIConfig(
			acct.Backend,
			acct.ProviderConfiguration.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)

		if c.Accounts[name].Tools != nil && c.Accounts[name].Tools.CircleCI != nil {
			ciConfig = ciConfig.merge(accountConfig, c.Accounts[name].Tools.CircleCI.Providers["aws"])
		} else {
			ciConfig = ciConfig.merge(accountConfig, awsProvider)
		}
	}

	for envName, env := range p.Envs {
		for cName, d := range env.Components {
			envConfig := d.CircleCI.generateCIConfig(
				d.Backend,
				d.ProviderConfiguration.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)

			if c.Envs[envName].Tools != nil && c.Envs[envName].Tools.CircleCI != nil {
				ciConfig = ciConfig.merge(envConfig, c.Envs[envName].Tools.CircleCI.Providers["aws"])
			} else {
				ciConfig = ciConfig.merge(envConfig, awsProvider)
			}
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
	}

	ciConfig = ciConfig.populateBuckets(numBuckets)
	return CircleCIConfig{
		CIConfig:        *ciConfig,
		SSHFingerprints: sshFingerprints,
	}
}

func (p *Plan) buildGitHubActionsConfig(c *v2.Config, foggVersion string) GitHubActionsCIConfig {
	var env map[string]string

	if c.Defaults.Tools != nil && c.Defaults.Tools.GitHubActionsCI != nil {
		env = c.Defaults.Tools.GitHubActionsCI.Env
	}

	ciConfig := &CIConfig{
		FoggVersion: foggVersion,
		Env:         env,
	}

	var awsProvider v2.CIProviderConfig

	if c.Defaults.Tools != nil && c.Defaults.Tools.GitHubActionsCI != nil {
		awsProvider = c.Defaults.Tools.GitHubActionsCI.Providers["aws"]
	}

	globalConfig := p.Global.GitHubActionsCI.generateCIConfig(
		p.Global.Backend,
		p.Global.ProviderConfiguration.AWS,
		"global",
		"terraform/global")

	if c.Global.Tools != nil && c.Global.Tools.GitHubActionsCI != nil {
		ciConfig = ciConfig.merge(globalConfig, c.Global.Tools.GitHubActionsCI.Providers["aws"])
	} else {
		ciConfig = ciConfig.merge(globalConfig, awsProvider)
	}

	for name, acct := range p.Accounts {
		accountConfig := acct.GitHubActionsCI.generateCIConfig(
			acct.Backend,
			acct.ProviderConfiguration.AWS,
			fmt.Sprintf("accounts/%s", name),
			fmt.Sprintf("terraform/accounts/%s", name),
		)

		if c.Accounts[name].Tools != nil && c.Accounts[name].Tools.GitHubActionsCI != nil {
			ciConfig = ciConfig.merge(accountConfig, c.Accounts[name].Tools.GitHubActionsCI.Providers["aws"])
		} else {
			ciConfig = ciConfig.merge(accountConfig, awsProvider)
		}
	}

	for envName, env := range p.Envs {
		for cName, d := range env.Components {
			envConfig := d.GitHubActionsCI.generateCIConfig(
				d.Backend,
				d.ProviderConfiguration.AWS,
				fmt.Sprintf("%s/%s", envName, cName),
				fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
			)

			if c.Envs[envName].Tools != nil && c.Envs[envName].Tools.GitHubActionsCI != nil {
				ciConfig = ciConfig.merge(envConfig, c.Envs[envName].Tools.GitHubActionsCI.Providers["aws"])
			} else {
				ciConfig = ciConfig.merge(envConfig, awsProvider)
			}
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
	}

	var (
		sshKeySecrets []string
		runsOn        []string
	)

	if c.Defaults.Tools != nil && c.Defaults.Tools.GitHubActionsCI != nil {
		sshKeySecrets = c.Defaults.Tools.GitHubActionsCI.SSHKeySecrets
	}

	if c.Defaults.Tools != nil &&
		c.Defaults.Tools.GitHubActionsCI != nil &&
		c.Defaults.Tools.GitHubActionsCI.RunsOn != nil {
		runsOn = *c.Defaults.Tools.GitHubActionsCI.RunsOn
	}

	ciConfig = ciConfig.populateBuckets(numBuckets)
	return GitHubActionsCIConfig{
		CIConfig:      *ciConfig,
		SSHKeySecrets: sshKeySecrets,
		RunsOn:        runsOn,
	}
}
