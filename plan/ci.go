package plan

import (
	"fmt"
	"path/filepath"
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
	TerraformVersion string
	SSHKeySecrets    []string
}

type AtlantisConfig struct {
	Enabled bool
	Envs    *map[string]Env
	RepoCfg *atlantis.RepoCfg
}

// alias v2.JavascriptPackageScope to jsScope for brevity
type jsScope = v2.JavascriptPackageScope
type TurboConfig struct {
	Enabled                 bool
	Version                 string
	RootName                string
	SCMBase                 string
	DevDependencies         map[string]string
	CdktfPackages           []string
	Workspaces              []vsCodeWorkspace
	Scopes                  map[string]jsScope
	CodeArtifactLoginScript string
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
		CIConfig:         *ciConfig,
		SSHKeySecrets:    sshKeySecrets,
		TerraformVersion: v2.ResolveRequiredString(v2.TerraformVersionGetter, c.Defaults.Common, c.Global.Common),
	}
}

const noCALoginRequired = "echo 'No CodeArtifact login required'"

type vsCodeWorkspace struct {
	Name string
	Path string
}

func (p *Plan) buildTurboRootConfig(c *v2.Config) *TurboConfig {
	turboConfig := &TurboConfig{
		Enabled:  false,
		SCMBase:  "main",
		RootName: "fogg-monorepo",
		DevDependencies: map[string]string{
			"turbo": "^2.4.0", // https://github.com/vercel/turborepo/releases
		},
		CodeArtifactLoginScript: noCALoginRequired,
		// Ensure vincenthsh/fogg's helper pkg scope
		Scopes: map[string]jsScope{
			"@vincenthsh": {
				Name:        "@vincenthsh",
				RegistryUrl: "https://npm.pkg.github.com",
			},
		},
	}

	if c.Turbo != nil {
		if c.Turbo.Enabled != nil {
			turboConfig.Enabled = *c.Turbo.Enabled
		}

		if c.Turbo.SCMBase != nil {
			turboConfig.SCMBase = *c.Turbo.SCMBase
		}

		if c.Turbo.Version != nil {
			turboConfig.Version = *c.Turbo.Version
		}

		if c.Turbo.RootName != nil {
			turboConfig.RootName = *c.Turbo.RootName
		}

		if c.Turbo.Scopes != nil {
			scopes, loginScript := parseJsScopes(&turboConfig.Scopes, c.Turbo.Scopes)
			turboConfig.Scopes = *scopes
			turboConfig.CodeArtifactLoginScript = loginScript
		}

		for _, dep := range c.Turbo.DevDependencies {
			turboConfig.DevDependencies[dep.Name] = dep.Version
		}

		pkgs := []string{}
		workspaces := []vsCodeWorkspace{}
		for module, modulePlan := range p.Modules {
			// applyModules implementation detail for pnpm-workspace.yaml
			path := fmt.Sprintf("%s/modules/%s", util.RootPath, module)
			pkgs = append(pkgs, path)
			kind := modulePlan.Kind.GetOrDefault()
			if kind == v2.ModuleKindCDKTF {
				workspaces = append(workspaces, vsCodeWorkspace{
					Name: fmt.Sprintf("module-%s", module),
					Path: path,
				})
			}
		}

		for env, envPlan := range p.Envs {
			for component, componentPlan := range envPlan.Components {
				// applyEnvs implementation detail for pnpm-workspace.yaml
				path := fmt.Sprintf("%s/envs/%s/%s", util.RootPath, env, component)
				pkgs = append(pkgs, path)
				kind := componentPlan.Kind.GetOrDefault()
				if kind == v2.ComponentKindCDKTF || kind == v2.ComponentKindTerraConstruct {
					workspaces = append(workspaces, vsCodeWorkspace{
						Name: fmt.Sprintf("env-%s-%s", env, component),
						Path: path,
					})
				}
			}
		}
		slices.Sort(pkgs)
		sort.Slice(workspaces, func(i, j int) bool {
			return workspaces[i].Name < workspaces[j].Name
		})
		turboConfig.CdktfPackages = pkgs
		turboConfig.Workspaces = workspaces
	}
	return turboConfig
}

// Merge scopes and detect CodeArtifact registries and generate optional ca:login script
func parseJsScopes(mergedScopes *map[string]jsScope, new []jsScope) (*map[string]v2.JavascriptPackageScope, string) {
	caRepos := map[string]util.CodeArtifactRepository{}
	for _, newScope := range new {
		(*mergedScopes)[newScope.Name] = newScope
		if util.IsCodeArtifactURL(newScope.RegistryUrl) {
			if _, exists := caRepos[newScope.Name]; !exists {
				caRepo, err := util.ParseRegistryUrl(newScope.Name, newScope.RegistryUrl)
				if err != nil {
					logrus.Warnf("Failed to parse CodeArtifact registry URL: %s", newScope.RegistryUrl)
					continue
				}
				caRepos[newScope.Name] = *caRepo
			}
		}
	}
	// Generate the npm ca:login script for CodeArtifact registries
	caLoginScript := noCALoginRequired
	if len(caRepos) > 0 {
		loginSteps := make([]string, 0, len(caRepos))
		for _, repo := range caRepos {
			loginSteps = append(loginSteps, repo.LoginCommand())
		}
		// sort for deterministic output
		sort.Strings(loginSteps)
		caLoginScript = strings.Join(loginSteps, ";")
	}

	return mergedScopes, caLoginScript
}

// buildAtlantisConfig must be build after Envs
func (p *Plan) buildAtlantisConfig(c *v2.Config) AtlantisConfig {
	enabled := false
	autoplanRemoteStates := false
	repoCfg := atlantis.RepoCfg{}
	if c.Defaults.Tools != nil && c.Defaults.Tools.Atlantis != nil {
		enabled = *c.Defaults.Tools.Atlantis.Enabled
		if c.Defaults.Tools.Atlantis.AutoplanRemoteStates != nil {
			autoplanRemoteStates = *c.Defaults.Tools.Atlantis.AutoplanRemoteStates
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
				whenModified := []string{"**/*.tf", "**/*.tf.json", "**/*.tfvars", "**/*.tfvars.json"}
				if d.AutoplanRelativeGlobs != nil {
					whenModified = append(whenModified, d.AutoplanRelativeGlobs...)
				}
				if d.AutoplanFiles != nil {
					for _, f := range d.AutoplanFiles {
						path := fmt.Sprintf("%s/envs/%s/%s", util.RootPath, envName, cName)
						relPath, _ := filepath.Rel(path, f)
						whenModified = append(whenModified, relPath)
					}
				}
				// if global autoplan remote states is disabled or
				// the component has no dependencies defined, explicitly ignore `remote-states.tf`
				if !autoplanRemoteStates || !d.HasDependsOn {
					whenModified = append(whenModified, "!remote-states.tf")
				}
				projects = append(projects, atlantis.Project{
					Name:              util.Ptr(fmt.Sprintf("%s_%s", envName, cName)),
					Dir:               util.Ptr(fmt.Sprintf("terraform/envs/%s/%s", envName, cName)),
					TerraformVersion:  &d.ComponentCommon.Common.TerraformVersion,
					Workspace:         util.Ptr(atlantis.DefaultWorkspace),
					ApplyRequirements: []string{atlantis.ApprovedRequirement},
					Autoplan: &atlantis.Autoplan{
						Enabled: util.Ptr(true),
						// Additional whenModified entries are added during module inspection in apply phase
						WhenModified: whenModified,
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
