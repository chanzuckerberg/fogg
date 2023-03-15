package plan

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/chanzuckerberg/go-misc/ptr"

	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

// Plan represents a set of actions to take
type Plan struct {
	Accounts        map[string]Account    `yaml:"account"`
	Envs            map[string]Env        `yaml:"envs"`
	Global          Component             `yaml:"global"`
	Modules         map[string]Module     `yaml:"modules"`
	TravisCI        TravisCIConfig        `yaml:"travis_ci"`
	CircleCI        CircleCIConfig        `yaml:"circleci_ci"`
	GitHubActionsCI GitHubActionsCIConfig `yaml:"github_actions_ci"`
	Atlantis        AtlantisConfig        `yaml:"atlantis"`
	Version         string                `yaml:"version"`
	TFE             *TFEConfig            `yaml:"tfe"`
}

// Common represents common fields
type Common struct {
	PathToRepoRoot   string `yaml:"path_to_repo_root"`
	TerraformVersion string `yaml:"terraform_version"`
}

// ComponentCommon represents common fields for components
type ComponentCommon struct {
	Common `yaml:",inline"`

	AccountBackends       map[string]Backend         `yaml:"account_backends"`
	Accounts              map[string]*json.Number    `yaml:"all_accounts"`
	Backend               Backend                    `yaml:"backend"`
	ComponentBackends     map[string]Backend         `yaml:"comonent_backends"`
	Env                   string                     ` yaml:"env"`
	ExtraVars             map[string]string          `yaml:"extra_vars"`
	Name                  string                     `yaml:"name"`
	Owner                 string                     `yaml:"owner"`
	Project               string                     `yaml:"project"`
	ProviderConfiguration ProviderConfiguration      `yaml:"providers_configuration"`
	ProviderVersions      map[string]ProviderVersion `yaml:"provider_versions"`

	TfLint TfLint `yaml:"tf_lint"`

	TravisCI        TravisCIComponent
	CircleCI        CircleCIComponent
	GitHubActionsCI GitHubActionsComponent
}

type TravisCIComponent struct {
	CIComponent
}

type CircleCIComponent struct {
	CIComponent

	SSHFingerprints []string
}

type GitHubActionsComponent struct {
	CIComponent
}

type CIComponent struct {
	Enabled     bool
	Buildevents bool

	AWSProfileName string
	AWSRoleName    string
	AWSRegion      string
	AWSAccountID   string
	Command        string
}

// generateCIConfig generates the config for ci
func (c CIComponent) generateCIConfig(
	backend Backend,
	provider *AWSProvider,
	projName string,
	projDir string) *CIConfig {
	if !c.Enabled {
		return nil
	}

	ciConfig := &CIConfig{
		AWSProfiles: ciAwsProfiles{},
		Enabled:     true,
		Buildevents: c.Buildevents,
	}

	ciConfig.projects = append(ciConfig.projects, CIProject{
		Name:    projName,
		Dir:     projDir,
		Command: c.Command,
	})

	if backend.S3 != nil && backend.S3.AccountID != nil && backend.S3.Profile != nil {
		p := *backend.S3.Profile
		ciConfig.AWSProfiles[p] = AWSRole{
			AccountID: *backend.S3.AccountID,
			RoleName:  c.AWSRoleName,
		}
	}

	if provider != nil {
		if provider.Profile != nil {
			ciConfig.AWSProfiles[*provider.Profile] = AWSRole{
				AccountID: provider.AccountID.String(),
				RoleName:  c.AWSRoleName,
			}
		}
	}
	return ciConfig
}

type ProviderConfiguration struct {
	Assert                 *AssertProvider     `yaml:"assert"`
	Auth0                  *Auth0Provider      `yaml:"auth0"`
	AWS                    *AWSProvider        `yaml:"aws"`
	AWSAdditionalProviders []AWSProvider       `yaml:"aws_regional_providers"`
	Bless                  *BlessProvider      `yaml:"bless"`
	Datadog                *DatadogProvider    `yaml:"datadog"`
	Github                 *GithubProvider     `yaml:"github"`
	Grafana                *GrafanaProvider    `yaml:"grafana"`
	Heroku                 *HerokuProvider     `yaml:"heroku"`
	Kubernetes             *KubernetesProvider `yaml:"kubernetes"`
	Okta                   *OktaProvider       `yaml:"okta"`
	Sentry                 *SentryProvider     `yaml:"sentry"`
	Snowflake              *SnowflakeProvider  `yaml:"snowflake"`
	Tfe                    *TfeProvider        `yaml:"tfe"`
	Sops                   *SopsProvider       `yaml:"sops"`
}

type ProviderVersion struct {
	Source  string  `yaml:"source"`
	Version *string `yaml:"version"`
}

var utilityProviders = map[string]ProviderVersion{
	"random": {
		Source:  "hashicorp/random",
		Version: ptr.String("~> 3.4"),
	},
	"archive": {
		Source:  "hashicorp/archive",
		Version: ptr.String("~> 2.0"),
	},
	"null": {
		Source:  "hashicorp/null",
		Version: ptr.String("~> 3.0"),
	},
	"local": {
		Source:  "hashicorp/local",
		Version: ptr.String("~> 2.0"),
	},
	"tls": {
		Source:  "hashicorp/tls",
		Version: ptr.String("~> 3.0"),
	},
	"assert": {
		Source:  "bwoznicki/assert",
		Version: ptr.String("~> 0.0.1"),
	},
	"okta-head": {
		Source:  "okta/okta",
		Version: ptr.String("~> 3.30"),
	},
}

// AWSProvider represents AWS provider configuration
type AWSProvider struct {
	AccountID json.Number `yaml:"account_id"`
	Alias     *string     `yaml:"alias"`
	Profile   *string     `yaml:"profile"`
	Region    string      `yaml:"region"`
	RoleArn   *string     `yaml:"role_arn"`
}

func (a *AWSProvider) String() string {
	accountID := a.AccountID.String()

	joined := util.JoinStrPointers(
		"-",
		&accountID,
		a.Alias,
		a.Profile,
		&a.Region,
		a.RoleArn,
	)
	if joined == nil {
		return ""
	}
	return *joined
}

// GithubProvider represents a configuration of a github provider
type GithubProvider struct {
	Organization string  `yaml:"organization"`
	BaseURL      *string `yaml:"base_url"`
}

type Auth0Provider struct {
	Domain string `yaml:"domain,omitempty"`
}

type AssertProvider struct {
	Version string `yaml:"version,omitempty"`
}

// SnowflakeProvider represents Snowflake DB provider configuration
type SnowflakeProvider struct {
	Account string `yaml:"account,omitempty"`
	Role    string `yaml:"role,omitempty"`
	Region  string `yaml:"region,omitempty"`
}

// OktaProvider represents Okta configuration
type OktaProvider struct {
	OrgName string  `yaml:"org_name,omitempty"`
	BaseURL *string `yaml:"base_url,omitempty"`
}

// BlessProvider represents Bless ssh provider configuration
type BlessProvider struct {
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `yaml:"aws_profile,omitempty"`
	AWSRegion         string   `yaml:"aws_region,omitempty"`
	RoleArn           *string  `yaml:"role_arn,omitempty"`
}

type HerokuProvider struct {
}

type DatadogProvider struct {
}

type SentryProvider struct {
	Enabled bool
	BaseURL *string `yaml:"base_url,omitempty"`
}

type TfeProvider struct {
	Enabled  bool    `yaml:"enabled,omitempty"`
	Hostname *string `yaml:"hostname,omitempty"`
}

type SopsProvider struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type KubernetesProvider struct {
}

type GrafanaProvider struct {
}

// BackendKind is a enum of backends we support
type BackendKind string

const (
	// BackendKindS3 is https://www.terraform.io/docs/backends/types/s3.html
	BackendKindS3     BackendKind = "s3"
	BackendKindRemote BackendKind = "remote"
)

// Backend represents a plan for configuring the terraform backend. Only one struct member can be
// non-nil at a time
type Backend struct {
	Kind   BackendKind    `yaml:"kind"`
	S3     *S3Backend     `yaml:"s3,omitempty"`
	Remote *RemoteBackend `yaml:"remote,omitempty"`
}

// S3Backend represents aws backend configuration
type S3Backend struct {
	AccountID   *string `yaml:"account_id,omitempty"`
	AccountName string  `yaml:"account_name"`
	Bucket      string  `yaml:"bucket"`
	DynamoTable *string `yaml:"dynamo_table"`
	KeyPath     string  `yaml:"key_path"`
	Profile     *string `yaml:"profile"`
	Region      string  `yaml:"region"`
	RoleArn     *string `yaml:"role_arn"`
}

// RemoteBackend represents a plan to configure a terraform remote backend
type RemoteBackend struct {
	HostName     string `yaml:"host_name"`
	Organization string `yaml:"organization"`
	Workspace    string `yaml:"workspace"`
}

// Module is a module
type Module struct {
	Common `yaml:",inline"`
}

// Account is an account
type Account struct {
	ComponentCommon `yaml:",inline"`

	Account string     `yaml:"account"`
	Global  *Component `yaml:"global"`
}

// Component is a component
type Component struct {
	ComponentCommon `yaml:",inline"`

	EKS          *v2.EKSConfig        `yaml:"eks,omitempty"`
	Kind         *v2.ComponentKind    `yaml:"kind,omitempty"`
	ModuleSource *string              `yaml:"module_source"`
	ModuleName   *string              `yaml:"module_name"`
	Variables    []string             `yaml:"variables"`
	Modules      []v2.ComponentModule `yaml:"modules"`
	Global       *Component           `yaml:"global"`
}

// Env is an env
type Env struct {
	Components map[string]Component `yaml:"components"`
	Env        string               `yaml:"env"`
	EKS        *v2.EKSConfig        `yaml:"eks"`
}

// TfLint contains a plan for running tflint
type TfLint struct {
	Enabled bool `yaml:"enabled"`
}

// Eval evaluates a config
func Eval(c *v2.Config) (*Plan, error) {
	if c == nil {
		return nil, errors.New("config is nil")
	}
	p := &Plan{}
	v, e := util.VersionString()
	if e != nil {
		return nil, errs.WrapInternal(e, "unable to parse fogg version")
	}
	p.Version = v

	var err error
	p.Global = p.buildGlobal(c)
	p.Accounts = p.buildAccounts(c)
	p.Envs, err = p.buildEnvs(c)
	if err != nil {
		return nil, err
	}

	p.Modules = p.buildModules(c)
	p.TravisCI = p.buildTravisCIConfig(c, v)
	p.CircleCI = p.buildCircleCIConfig(c, v)
	p.GitHubActionsCI = p.buildGitHubActionsConfig(c, v)
	p.Atlantis = p.buildAtlantisConfig(c, v)
	p.TFE, err = p.buildTFE(c)
	if err != nil {
		return p, err
	}
	return p, nil
}

// Print prints a plan
func Print(p *Plan) error {
	out, err := yaml.Marshal(p)
	if err != nil {
		return errs.WrapInternal(err, "yaml: could not marshal")
	}
	fmt.Print(string(out))
	return nil
}

func parseGithubOrgRepoFromGit(path string) (string, string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to open the git repo")
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return "", "", errors.Wrap(err, "unable to list the remotes of the repo")
	}
	for _, remote := range remotes {
		if remote.Config().Name != "origin" {
			continue
		}

		for _, u := range remote.Config().URLs {
			remoteSplit := strings.Split(u, ":")
			if len(remoteSplit) != 2 {
				return "", "", errors.Errorf("unexpected syntax in git remote URL %s", u)
			}
			baseName := filepath.Base(remoteSplit[1])
			return filepath.Dir(remoteSplit[1]), strings.TrimSuffix(baseName, filepath.Ext(remoteSplit[1])), nil
		}
	}

	return "", "", errors.New("unable to find a valid origin remote URL")
}

func (p *Plan) buildTFE(c *v2.Config) (*TFEConfig, error) {
	if c.TFE == nil {
		return nil, nil
	}

	tfeConfig := &TFEConfig{
		ReadTeams:                      []string{},
		Branch:                         "main",
		GithubOrg:                      "",
		GithubRepo:                     "",
		TFEOrg:                         c.TFE.TFEOrg,
		SSHKeyName:                     "fogg-ssh-key",
		ExcludedGithubRequiredChecks:   []string{},
		AdditionalGithubRequiredChecks: []string{},
	}
	tfeConfig.ComponentCommon = resolveComponentCommon(c.Defaults.Common, c.Global.Common, c.TFE.Common)
	tfeConfig.ModuleSource = c.TFE.ModuleSource
	tfeConfig.ModuleName = c.TFE.ModuleName

	if tfeConfig.ComponentCommon.Backend.Kind == BackendKindS3 {
		tfeConfig.ComponentCommon.Backend.S3.KeyPath = fmt.Sprintf("terraform/%s/%s.tfstate", "tfe", "tfe")
	} else if tfeConfig.ComponentCommon.Backend.Kind == BackendKindRemote {
		tfeConfig.ComponentCommon.Backend.Remote.Workspace = "tfe"
	} else {
		panic(fmt.Sprintf("Invalid backend kind of %s", tfeConfig.ComponentCommon.Backend.Kind))
	}

	if c.TFE.ReadTeams != nil {
		tfeConfig.ReadTeams = *c.TFE.ReadTeams
	}
	if c.TFE.Branch != nil {
		tfeConfig.Branch = *c.TFE.Branch
	}
	if c.TFE.GithubOrg != nil {
		tfeConfig.GithubOrg = *c.TFE.GithubOrg
	}
	if c.TFE.GithubRepo != nil {
		tfeConfig.GithubRepo = *c.TFE.GithubRepo
	}
	if c.TFE.SSHKeyName != nil {
		tfeConfig.SSHKeyName = *c.TFE.SSHKeyName
	}
	if c.TFE.ExcludedGithubRequiredChecks != nil {
		tfeConfig.ExcludedGithubRequiredChecks = *c.TFE.ExcludedGithubRequiredChecks
	}
	if c.TFE.AdditionalGithubRequiredChecks != nil {
		tfeConfig.AdditionalGithubRequiredChecks = *c.TFE.AdditionalGithubRequiredChecks
	}

	if tfeConfig.GithubOrg == "" || tfeConfig.GithubRepo == "" {
		org, repo, err := parseGithubOrgRepoFromGit(".")
		if err != nil {
			return nil, err
		}
		tfeConfig.GithubOrg = org
		tfeConfig.GithubRepo = repo
	}

	return tfeConfig, nil
}

func (p *Plan) buildAccounts(c *v2.Config) map[string]Account {
	defaults := c.Defaults

	accountPlans := make(map[string]Account, len(c.Accounts))
	for name, acct := range c.Accounts {
		accountPlan := Account{}

		accountPlan.ComponentCommon = resolveComponentCommon(defaults.Common, acct.Common)
		accountPlan.Name = name
		accountPlan.Account = name // for backwards compat
		accountPlan.Env = "accounts"

		if accountPlan.ComponentCommon.Backend.Kind == BackendKindS3 {
			accountPlan.ComponentCommon.Backend.S3.KeyPath = fmt.Sprintf("terraform/%s/accounts/%s.tfstate", accountPlan.ComponentCommon.Project, name)
		} else if accountPlan.ComponentCommon.Backend.Kind == BackendKindRemote {
			accountPlan.ComponentCommon.Backend.Remote.Workspace = fmt.Sprintf("accounts-%s", name)
		} else {
			panic(fmt.Sprintf("Invalid backend kind of %s", accountPlan.ComponentCommon.Backend.Kind))
		}

		accountPlan.Accounts = resolveAccounts(c.Accounts) //FIXME this needs to run as a second phase, not directly from the config
		accountPlan.PathToRepoRoot = "../../../"
		accountPlan.Global = &p.Global
		accountPlans[name] = accountPlan
	}

	accountBackends := make(map[string]Backend)
	for acctName, acct := range accountPlans {
		accountBackends[acctName] = acct.Backend
	}

	for name, acct := range c.Accounts {
		accountRemoteStates := v2.ResolveOptionalStringSlice(v2.DependsOnAccountsGetter, defaults.Common, acct.Common)
		a := accountPlans[name]
		filtered := map[string]Backend{}

		if accountRemoteStates != nil {
			for k, v := range accountBackends {
				if util.SliceContainsString(accountRemoteStates, k) {
					filtered[k] = v
				}
			}
		} else {
			filtered = accountBackends
		}
		a.AccountBackends = filtered
		accountPlans[name] = a
	}

	return accountPlans
}

func (p *Plan) buildModules(c *v2.Config) map[string]Module {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.PathToRepoRoot = "../../../"
		modulePlan.TerraformVersion = *v2.ResolveModuleTerraformVersion(c.Defaults, conf)
		modulePlans[name] = modulePlan
	}
	return modulePlans
}

func newEnvPlan() Env {
	ep := Env{}
	ep.Components = make(map[string]Component)
	return ep
}

func (p *Plan) buildGlobal(conf *v2.Config) Component {
	// Global just uses defaults because that's the way sicc worked. We should make it directly configurable.
	componentPlan := Component{}
	componentPlan.Accounts = resolveAccounts(conf.Accounts)
	defaults := conf.Defaults
	global := conf.Global

	componentPlan.ComponentCommon = resolveComponentCommon(defaults.Common, global.Common)

	if componentPlan.ComponentCommon.Backend.Kind == BackendKindS3 {
		componentPlan.ComponentCommon.Backend.S3.KeyPath = fmt.Sprintf("terraform/%s/global.tfstate", componentPlan.ComponentCommon.Project)
	} else if componentPlan.ComponentCommon.Backend.Kind == BackendKindRemote {
		componentPlan.ComponentCommon.Backend.Remote.Workspace = "global"
	}

	componentPlan.Name = "global"
	componentPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, global.ExtraVars)
	componentPlan.PathToRepoRoot = "../../"

	return componentPlan
}

// buildEnvs must be build after accounts
func (p *Plan) buildEnvs(conf *v2.Config) (map[string]Env, error) {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults

	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()
		envPlan.Env = envName

		for componentName, componentConf := range conf.Envs[envName].Components {
			componentPlan := Component{
				Kind: componentConf.Kind,
			}
			// reference accounts for remote state
			if _, dupe := p.Accounts[componentName]; dupe {
				return nil, errs.WrapUser(fmt.Errorf("Component %s can't have same name as account", componentName), "Invalid component name")
			}

			componentPlan.ComponentCommon = resolveComponentCommon(defaults.Common, envConf.Common, componentConf.Common)
			accountRemoteStates := v2.ResolveOptionalStringSlice(v2.DependsOnAccountsGetter, defaults.Common, envConf.Common, componentConf.Common)
			accountBackends := map[string]Backend{}
			for k, v := range p.Accounts {
				if accountRemoteStates == nil || util.SliceContainsString(accountRemoteStates, k) {
					accountBackends[k] = v.Backend
				}
			}
			componentPlan.AccountBackends = accountBackends

			componentPlan.Accounts = resolveAccounts(conf.Accounts)

			if componentPlan.ComponentCommon.Backend.Kind == BackendKindS3 {
				componentPlan.ComponentCommon.Backend.S3.KeyPath = fmt.Sprintf("terraform/%s/envs/%s/components/%s.tfstate", componentPlan.ComponentCommon.Project, envName, componentName)
			} else if componentPlan.ComponentCommon.Backend.Kind == BackendKindRemote {
				componentPlan.ComponentCommon.Backend.Remote.Workspace = fmt.Sprintf("%s-%s", envName, componentName)
			} else {
				panic(fmt.Sprintf("Invalid backend kind of %s", componentPlan.ComponentCommon.Backend.Kind))
			}

			componentPlan.Env = envName
			componentPlan.Name = componentName
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.ModuleName = componentConf.ModuleName
			componentPlan.Variables = componentConf.Variables
			componentPlan.Modules = componentConf.Modules
			componentPlan.PathToRepoRoot = "../../../../"

			if !envConf.NoGlobal {
				componentPlan.Global = &p.Global
			}

			envPlan.Components[componentName] = componentPlan
		}

		componentBackends := make(map[string]Backend)

		for componentName, component := range envPlan.Components {
			// FIXME (el): get rid of non-terraform component kinds
			if component.Kind.GetOrDefault() != v2.ComponentKindTerraform {
				continue
			}

			componentBackends[componentName] = component.Backend
		}

		for name, componentConf := range conf.Envs[envName].Components {
			componentRemoteStates := v2.ResolveOptionalStringSlice(v2.DependsOnComponentsGetter, defaults.Common, envConf.Common, componentConf.Common)
			c := envPlan.Components[name]
			filtered := map[string]Backend{}

			if componentRemoteStates != nil {
				for k, v := range componentBackends {
					if util.SliceContainsString(componentRemoteStates, k) {
						filtered[k] = v
					}
				}
			} else {
				filtered = componentBackends
			}

			c.ComponentBackends = filtered
			envPlan.Components[name] = c
		}

		envPlans[envName] = envPlan
	}
	return envPlans, nil
}

func resolveAWSProvider(commons ...v2.Common) (plan *AWSProvider, providers []AWSProvider, version *string) {
	awsConfig := v2.ResolveAWSProvider(commons...)
	var roleArn *string
	// nothing to do
	if awsConfig == nil {
		return
	}

	// set the version
	version = awsConfig.Version

	// configure the main provider
	if awsConfig.Role != nil {
		tmp := fmt.Sprintf("arn:aws:iam::%s:role/%s", *awsConfig.AccountID, *awsConfig.Role)
		roleArn = &tmp
	}
	plan = &AWSProvider{
		AccountID: *awsConfig.AccountID,
		Profile:   awsConfig.Profile,
		Region:    *awsConfig.Region,
		RoleArn:   roleArn,
	}

	// grab all aliased regions
	for _, r := range awsConfig.AdditionalRegions {
		// we have to take a reference here otherwise it gets overwritten by the loop
		region := r
		providers = append(providers,
			AWSProvider{
				AccountID: *awsConfig.AccountID,
				Alias:     &region,
				Profile:   awsConfig.Profile,
				Region:    region,
				RoleArn:   roleArn,
			})
	}

	//HACK HACK(el): this is horrible: grab all extra accounts and configure aliased providers for them
	for ap, aws := range awsConfig.AdditionalProviders {
		aliasPrefix := ap
		// HACK(el): we create this pseudo v2.Common for each additional provider and do the inheritance
		pseudoCommon := v2.Common{
			Providers: &v2.Providers{
				AWS: aws,
			},
		}
		// set so not nil
		pseudoCommon.Providers.AWS.AdditionalProviders = map[string]*v2.AWSProvider{}

		extraCommons := append(commons[:], pseudoCommon)
		extraPlan, extraProviders, _ := resolveAWSProvider(extraCommons...)
		if extraPlan != nil {
			extraPlan.Alias = util.JoinStrPointers("-", &aliasPrefix, extraPlan.Alias)
			providers = append(providers, *extraPlan)
		}

		for _, eP := range extraProviders {
			extraProvider := eP
			extraProvider.Alias = util.JoinStrPointers("-", &aliasPrefix, extraProvider.Alias)
			providers = append(providers, extraProvider)
		}
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].String() > providers[j].String()
	})
	return
}

func resolveComponentCommon(commons ...v2.Common) ComponentCommon {
	providerVersions := copyMap(utilityProviders)
	awsPlan, additionalProviders, awsVersion := resolveAWSProvider(commons...)
	if awsVersion != nil {
		providerVersions["aws"] = ProviderVersion{
			Source:  "hashicorp/aws",
			Version: awsVersion,
		}
	}

	var auth0Plan *Auth0Provider
	auth0Config := v2.ResolveAuth0Provider(commons...)
	if auth0Config != nil {
		auth0Plan = &Auth0Provider{
			Domain: *auth0Config.Domain,
		}

		defaultSource := "alexkappa/auth0"
		if auth0Config.Source == nil {
			auth0Config.Source = &defaultSource
		}
		providerVersions["auth0"] = ProviderVersion{
			Source:  *auth0Config.Source,
			Version: auth0Config.Version,
		}
	}

	var assertPlan *AssertProvider
	assertConfig := v2.ResolveAssertProvider(commons...)
	if assertConfig != nil {
		assertPlan = &AssertProvider{
			Version: *assertConfig.Version,
		}

		providerVersions["assert"] = ProviderVersion{
			Source:  "bwoznicki/assert",
			Version: assertConfig.Version,
		}
	}

	var githubPlan *GithubProvider
	githubConfig := v2.ResolveGithubProvider(commons...)

	if githubConfig != nil {
		githubPlan = &GithubProvider{
			Organization: *githubConfig.Organization,
			BaseURL:      githubConfig.BaseURL,
		}

		providerVersions["github"] = ProviderVersion{
			Source:  "integrations/github",
			Version: githubConfig.Version,
		}
	}

	var snowflakePlan *SnowflakeProvider
	snowflakeConfig := v2.ResolveSnowflakeProvider(commons...)
	if snowflakeConfig != nil {
		snowflakePlan = &SnowflakeProvider{
			Account: *snowflakeConfig.Account,
			Role:    *snowflakeConfig.Role,
			Region:  *snowflakeConfig.Region,
		}

		providerVersions["snowflake"] = ProviderVersion{
			Source:  "Snowflake-Labs/snowflake",
			Version: snowflakeConfig.Version,
		}
	}

	var oktaPlan *OktaProvider
	oktaConfig := v2.ResolveOktaProvider(commons...)
	if oktaConfig != nil {
		oktaPlan = &OktaProvider{
			OrgName: *oktaConfig.OrgName,
			BaseURL: oktaConfig.BaseURL,
		}

		var registryNamespace string

		if oktaConfig.RegistryNamespace != nil && *oktaConfig.RegistryNamespace != "" {
			registryNamespace = *oktaConfig.RegistryNamespace
		} else {
			registryNamespace = "oktadeveloper"
		}
		providerVersions["okta"] = ProviderVersion{
			Source:  fmt.Sprintf("%s/okta", registryNamespace),
			Version: oktaConfig.Version,
		}
	}

	var blessPlan *BlessProvider
	blessConfig := v2.ResolveBlessProvider(commons...)
	if blessConfig != nil && (blessConfig.AWSProfile != nil || blessConfig.RoleArn != nil) && blessConfig.AWSRegion != nil {
		blessPlan = &BlessProvider{
			AWSProfile:        blessConfig.AWSProfile,
			AWSRegion:         *blessConfig.AWSRegion,
			AdditionalRegions: blessConfig.AdditionalRegions,
			RoleArn:           blessConfig.RoleArn,
		}

		providerVersions["bless"] = ProviderVersion{
			Source:  "chanzuckerberg/bless",
			Version: blessConfig.Version,
		}
	}

	var herokuPlan *HerokuProvider
	herokuConfig := v2.ResolveHerokuProvider(commons...)
	if herokuConfig != nil {
		herokuPlan = &HerokuProvider{}

		providerVersions["heroku"] = ProviderVersion{
			Source:  "heroku/heroku",
			Version: herokuConfig.Version,
		}
	}

	var datadogPlan *DatadogProvider
	datadogConfig := v2.ResolveDatadogProvider(commons...)
	if datadogConfig != nil {
		datadogPlan = &DatadogProvider{}

		providerVersions["datadog"] = ProviderVersion{
			Source:  "datadog/datadog",
			Version: datadogConfig.Version,
		}
	}

	pagerdutyConfig := v2.ResolvePagerdutyProvider(commons...)
	if pagerdutyConfig != nil {
		providerVersions["pagerduty"] = ProviderVersion{
			Source:  "pagerduty/pagerduty",
			Version: pagerdutyConfig.Version,
		}
	}

	opsGenieConfig := v2.ResolveOpsGenieProvider(commons...)
	if opsGenieConfig != nil {
		providerVersions["opsgenie"] = ProviderVersion{
			Source:  "opsgenie/opsgenie",
			Version: opsGenieConfig.Version,
		}
	}

	databricksConfig := v2.ResolveDatabricksProvider(commons...)
	if databricksConfig != nil {
		providerVersions["databricks"] = ProviderVersion{
			Source:  "databricks/databricks",
			Version: databricksConfig.Version,
		}
	}

	var sentryPlan *SentryProvider
	sentryConfig := v2.ResolveSentryProvider(commons...)
	if sentryConfig != nil {
		sentryPlan = &SentryProvider{
			Enabled: true,
			BaseURL: sentryConfig.BaseURL,
		}

		providerVersions["sentry"] = ProviderVersion{
			Source:  "jianyuan/sentry",
			Version: sentryConfig.Version,
		}
	}

	var tfePlan *TfeProvider

	tfeConfig := v2.ResolveTfeProvider(commons...)
	if tfeConfig.Enabled != nil && *tfeConfig.Enabled {
		tfePlan = &TfeProvider{
			Enabled:  true,
			Hostname: tfeConfig.Hostname,
		}

		providerVersions["tfe"] = ProviderVersion{
			Source:  "hashicorp/tfe",
			Version: tfeConfig.Version,
		}
	}

	var sopsPlan *SopsProvider

	sopsConfig := v2.ResolveSopsProvider(commons...)
	if sopsConfig.Enabled != nil && *sopsConfig.Enabled {
		sopsPlan = &SopsProvider{
			Enabled: true,
		}

		providerVersions["sops"] = ProviderVersion{
			Source:  "carlpett/sops",
			Version: sopsConfig.Version,
		}
	}

	var k8sPlan *KubernetesProvider

	k8sConfig := v2.ResolveKubernetesProvider(commons...)
	if k8sConfig.Enabled != nil && *k8sConfig.Enabled {
		k8sPlan = &KubernetesProvider{}

		providerVersions["kubernetes"] = ProviderVersion{
			Source:  "hashicorp/kubernetes",
			Version: k8sConfig.Version,
		}
	}

	var grafanaPlan *GrafanaProvider

	grafanaConfig := v2.ResolveGrafanaProvider(commons...)
	if grafanaConfig.Enabled != nil && *grafanaConfig.Enabled {
		grafanaPlan = &GrafanaProvider{}

		providerVersions["grafana"] = ProviderVersion{
			Source:  "grafana/grafana",
			Version: grafanaConfig.Version,
		}
	}

	tflintConfig := v2.ResolveTfLint(commons...)

	tfLintPlan := TfLint{
		Enabled: *tflintConfig.Enabled,
	}

	travisConfig := v2.ResolveTravis(commons...)
	travisPlan := TravisCIComponent{
		CIComponent: CIComponent{
			Enabled:     *travisConfig.Enabled,
			Buildevents: *travisConfig.Buildevents,
		},
	}
	if travisPlan.Enabled {
		travisPlan.AWSRoleName = *travisConfig.AWSIAMRoleName
		travisPlan.Command = *travisConfig.Command
	}

	githubActionsConfig := v2.ResolveGitHubActionsCI(commons...)
	githubActionsPlan := GitHubActionsComponent{
		CIComponent: CIComponent{
			Enabled:     *githubActionsConfig.Enabled,
			Buildevents: *githubActionsConfig.Buildevents,
		},
	}
	if githubActionsPlan.Enabled {
		if githubActionsConfig.AWSIAMRoleName != nil {
			githubActionsPlan.AWSRoleName = *githubActionsConfig.AWSIAMRoleName
			githubActionsPlan.AWSRegion = *githubActionsConfig.AWSRegion
		}
		githubActionsPlan.Command = *githubActionsConfig.Command
	}

	circleConfig := v2.ResolveCircleCI(commons...)
	circlePlan := CircleCIComponent{
		CIComponent: CIComponent{
			Enabled:     *circleConfig.Enabled,
			Buildevents: *circleConfig.Buildevents,
		},
	}
	if circlePlan.Enabled {
		circlePlan.AWSRoleName = *circleConfig.AWSIAMRoleName
		circlePlan.Command = *circleConfig.Command
	}

	if travisPlan.Enabled && circlePlan.Enabled {
		logrus.Warn("Detected both travisCI and circleCI are enabled, is this intentional?")
	}

	project := v2.ResolveRequiredString(v2.ProjectGetter, commons...)

	backendConf := v2.ResolveBackend(commons...)
	var backend Backend

	if backendConf != nil {
		if *backendConf.Kind == "s3" {
			var roleArn *string
			if backendConf.Role != nil {
				// we know from our validations that if role is set, then account id must also be set
				tmp := fmt.Sprintf("arn:aws:iam::%s:role/%s", *backendConf.AccountID, *backendConf.Role)
				roleArn = &tmp
			}
			backend = Backend{
				Kind: BackendKindS3,
				S3: &S3Backend{
					Region:  *backendConf.Region,
					Profile: backendConf.Profile,
					Bucket:  *backendConf.Bucket,
					RoleArn: roleArn,

					AccountID:   backendConf.AccountID,
					DynamoTable: backendConf.DynamoTable,
				},
			}
		} else if *backendConf.Kind == "remote" {
			backend = Backend{
				Kind: BackendKindRemote,
				Remote: &RemoteBackend{
					HostName:     *backendConf.HostName,
					Organization: *backendConf.Organization,
				},
			}
		}
	}

	return ComponentCommon{
		Backend: backend,
		ProviderConfiguration: ProviderConfiguration{
			Assert:                 assertPlan,
			Auth0:                  auth0Plan,
			AWS:                    awsPlan,
			AWSAdditionalProviders: additionalProviders,
			Bless:                  blessPlan,
			Datadog:                datadogPlan,
			Github:                 githubPlan,
			Grafana:                grafanaPlan,
			Heroku:                 herokuPlan,
			Kubernetes:             k8sPlan,
			Okta:                   oktaPlan,
			Sentry:                 sentryPlan,
			Snowflake:              snowflakePlan,
			Tfe:                    tfePlan,
			Sops:                   sopsPlan,
		},
		ProviderVersions: providerVersions,
		TfLint:           tfLintPlan,
		ExtraVars:        v2.ResolveStringMap(v2.ExtraVarsGetter, commons...),
		Owner:            v2.ResolveRequiredString(v2.OwnerGetter, commons...),
		Project:          project,
		Common:           Common{TerraformVersion: v2.ResolveRequiredString(v2.TerraformVersionGetter, commons...)},
		TravisCI:         travisPlan,
		CircleCI:         circlePlan,
		GitHubActionsCI:  githubActionsPlan,
	}
}

func resolveExtraVars(vars ...map[string]string) map[string]string {
	resolved := map[string]string{}

	for _, m := range vars {
		for k, v := range m {
			resolved[k] = v
		}
	}
	return resolved
}

func resolveAccounts(accounts map[string]v2.Account) map[string]*json.Number {
	a := make(map[string]*json.Number)
	for name, account := range accounts {
		if account.Providers != nil && account.Providers.AWS != nil && account.Providers.AWS.AccountID != nil {
			a[name] = account.Providers.AWS.AccountID
		}
	}
	return a
}

func copyMap(in map[string]ProviderVersion) map[string]ProviderVersion {
	out := map[string]ProviderVersion{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
