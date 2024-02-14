package plan

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	atlantis "github.com/runatlantis/atlantis/server/core/config/raw"
	"github.com/sirupsen/logrus"
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

	DefaultAWSIAMRoleName string
	DefaultAWSRegion      string

	TestBuckets [][]CIProject
	AWSProfiles ciAwsProfiles
	Buildevents bool

	PreCommit PreCommitConfig
}

type CircleCIConfig struct {
	CIConfig
	SSHFingerprints []string
}

type GitHubActionsCIConfig struct {
	CIConfig
	SSHKeySecrets []string
}

type AtlantisConfig struct {
	Enabled bool
	Envs    *map[string]Env
	RepoCfg *atlantis.RepoCfg
}

type PreCommitConfig struct {
	Enabled            bool
	Requirements       []string
	Config             *v2.PreCommitConfig
	GitHubActionSteps  []v2.GitHubActionStep
	HooksSkippedInMake []string
	ExtraArgs          []string
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
		PreCommit:   p.buildGithubActionsPreCommitConfig(c, foggVersion),
	}

	if c.Defaults.Tools != nil && c.Defaults.Tools.GitHubActionsCI != nil &&
		c.Defaults.Tools.GitHubActionsCI.AWSIAMRoleName != nil && c.Defaults.Tools.GitHubActionsCI.AWSRegion != nil {
		ciConfig.DefaultAWSIAMRoleName = *c.Defaults.Tools.GitHubActionsCI.AWSIAMRoleName
		ciConfig.DefaultAWSRegion = *c.Defaults.Tools.GitHubActionsCI.AWSRegion
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

	var sshKeySecrets []string

	if c.Defaults.Tools != nil && c.Defaults.Tools.GitHubActionsCI != nil {
		sshKeySecrets = c.Defaults.Tools.GitHubActionsCI.SSHKeySecrets
	}

	ciConfig = ciConfig.populateBuckets(numBuckets)
	return GitHubActionsCIConfig{
		CIConfig:      *ciConfig,
		SSHKeySecrets: sshKeySecrets,
	}
}

// buildAtlantisConfig must be build after Envs
func (p *Plan) buildAtlantisConfig(c *v2.Config, foggVersion string) AtlantisConfig {
	enabled := false
	repoCfg := atlantis.RepoCfg{}
	if c.Defaults.Tools != nil && c.Defaults.Tools.Atlantis != nil {
		enabled = *c.Defaults.Tools.Atlantis.Enabled
		modulePrefixes := c.Defaults.Tools.Atlantis.ModulePrefixes
		if len(modulePrefixes) == 0 {
			modulePrefixes = append(modulePrefixes, "terraform/modules/")
		}
		repoCfg = c.Defaults.Tools.Atlantis.RepoCfg
		projects := []atlantis.Project{}
		for envName, env := range p.Envs {
			for cName, d := range env.Components {
				uniqueModuleSources := []string{}
				if d.ModuleSource != nil {
					uniqueModuleSources = append(uniqueModuleSources, *d.ModuleSource)
				}
				for _, m := range d.Modules {
					if !slices.Contains(uniqueModuleSources, *m.Source) {
						uniqueModuleSources = append(uniqueModuleSources, *m.Source)
					}
				}

				projects = append(projects, atlantis.Project{
					Name:              util.Ptr(fmt.Sprintf("%s_%s", envName, cName)),
					Dir:               util.Ptr(fmt.Sprintf("terraform/envs/%s/%s", envName, cName)),
					TerraformVersion:  &d.ComponentCommon.Common.TerraformVersion,
					Workspace:         util.Ptr(atlantis.DefaultWorkspace),
					ApplyRequirements: []string{atlantis.ApprovedRequirement},
					Autoplan: &atlantis.Autoplan{
						Enabled:      util.Ptr(true),
						WhenModified: generateWhenModified(uniqueModuleSources, d.PathToRepoRoot, modulePrefixes),
					},
				})
			}
		}

		// sort projects by name
		sort.Slice(projects, func(i, j int) bool {
			return *projects[i].Name < *projects[j].Name
		})
		repoCfg.Projects = projects
	}
	return AtlantisConfig{
		Enabled: enabled,
		Envs:    &p.Envs,
		RepoCfg: &repoCfg,
	}
}

func generateWhenModified(moduleSources []string, pathToRepoRoot string, modulePrefixes []string) []string {
	whenModified := []string{
		"*.tf",
		"!remote-states.tf",
	}
	for _, moduleSource := range moduleSources {
		if startsWithPrefix(moduleSource, modulePrefixes) {
			modulePath := pathToRepoRoot + moduleSource
			whenModified = append(whenModified,
				fmt.Sprintf(
					"%s/**/*.tf",
					modulePath,
				), fmt.Sprintf(
					"%s/**/*.tf.json",
					modulePath,
				),
			)
		} else {
			logrus.Debugf("atlantis: moduleSource %q is not part of module_prefix list: %q", moduleSource, modulePrefixes)
		}
	}
	return whenModified
}

// startsWithPrefix checks if the given string s starts with any of the prefixes in the array.
func startsWithPrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func (p *Plan) buildGithubActionsPreCommitConfig(c *v2.Config, foggVersion string) PreCommitConfig {
	// defaults
	enabled := false
	config := v2.PreCommitConfig{}
	preCommitVersion := "3.4.0"
	requirements := []string{}
	steps := []v2.GitHubActionStep{}
	hooksSkippedInMake := []string{}
	extraArgs := []string{}

	if c.Defaults.Tools != nil &&
		c.Defaults.Tools.GitHubActionsCI != nil &&
		c.Defaults.Tools.GitHubActionsCI.PreCommit != nil {
		setup := c.Defaults.Tools.GitHubActionsCI.PreCommit
		enabled = setup.Enabled
		config = *setup.Config

		if setup.Version != nil {
			preCommitVersion = *setup.Version
		}

		if setup.GitHubActionSteps != nil && len(setup.GitHubActionSteps) > 0 {
			steps = append(steps, setup.GitHubActionSteps...)
		}

		if setup.ExtraArgs != nil && len(setup.ExtraArgs) > 0 {
			extraArgs = append(extraArgs, setup.ExtraArgs...)
		}

		if len(config.Repos) > 0 {
			for ri, r := range config.Repos {
				if r.Hooks == nil || len(r.Hooks) == 0 {
					continue
				}
				for hi, h := range r.Hooks {
					if h.SkipInMake != nil && *h.SkipInMake {
						hooksSkippedInMake = append(hooksSkippedInMake, h.ID)
					}
					// drop from yaml output to avoid pre-commit invalid field warnings
					config.Repos[ri].Hooks[hi].SkipInMake = nil
				}
			}
		}

		// ensure pre-commit is in requirements.txt
		requirements = append(requirements,
			fmt.Sprintf("pre-commit==%s", preCommitVersion),
		)
		// add any other dependencies (for GH setup-python Action with pip cache)
		for pkg, version := range setup.PipCache {
			requirements = append(requirements,
				fmt.Sprintf("%s==%s", pkg, version),
			)
		}
	}

	return PreCommitConfig{
		Enabled:            enabled,
		Requirements:       requirements,
		Config:             &config,
		GitHubActionSteps:  steps,
		HooksSkippedInMake: hooksSkippedInMake,
		ExtraArgs:          extraArgs,
	}
}
