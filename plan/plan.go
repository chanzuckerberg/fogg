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
	Version         string                `yaml:"version"`
	TFE             *TFEConfig            `yaml:"tfe"`
}

// Common represents common fields
type Common struct {
	PathToRepoRoot   string `yaml:"path_to_repo_root"`
	TerraformVersion string `yaml:"terraform_version"`
}

type ExtraTemplate struct {
	Overwrite bool
	Content   string
}

// ComponentCommon represents common fields for components
type ComponentCommon struct {
	Common `yaml:",inline"`

	AccountBackends          map[string]Backend         `yaml:"account_backends"`
	Accounts                 map[string]*json.Number    `yaml:"all_accounts"`
	Backend                  Backend                    `yaml:"backend"`
	ComponentBackends        map[string]Backend         `yaml:"comonent_backends"`
	Env                      string                     ` yaml:"env"`
	ExtraVars                map[string]string          `yaml:"extra_vars"`
	Name                     string                     `yaml:"name"`
	Owner                    string                     `yaml:"owner"`
	Project                  string                     `yaml:"project"`
	ProviderConfiguration    ProviderConfiguration      `yaml:"providers_configuration"`
	ProviderVersions         map[string]ProviderVersion `yaml:"provider_versions"`
	NeedsAWSAccountsVariable bool                       `yaml:"needs_aws_accounts_variable"`
	ExtraTemplates           map[string]ExtraTemplate   `yaml:"extra_templates"`

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
	Helm                   *HelmProvider       `yaml:"helm"`
	Kubectl                *KubectlProvider    `yaml:"kubectl"`
	Okta                   *OktaProvider       `yaml:"okta"`
	Sentry                 *SentryProvider     `yaml:"sentry"`
	Snowflake              *SnowflakeProvider  `yaml:"snowflake"`
	Tfe                    *TfeProvider        `yaml:"tfe"`
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
		Version: ptr.String("3.1.1"),
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
		Version: ptr.String("> 3.30"),
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
	CommonProvider `yaml:",inline"`
	Organization   string  `yaml:"organization"`
	BaseURL        *string `yaml:"base_url"`
}

type Auth0Provider struct {
	CommonProvider `yaml:",inline"`
	Domain         string `yaml:"domain,omitempty"`
}

type AssertProvider struct {
	CommonProvider `yaml:",inline"`
}

type CommonProvider struct {
	CustomProvider bool   `yaml:"custom_provider,omitempty"`
	Enabled        bool   `yaml:"enabled,omitempty"`
	Version        string `yaml:"version,omitempty"`
}

// SnowflakeProvider represents Snowflake DB provider configuration
type SnowflakeProvider struct {
	CommonProvider `yaml:",inline"`
	Account        string `yaml:"account,omitempty"`
	Role           string `yaml:"role,omitempty"`
	Region         string `yaml:"region,omitempty"`
}

// OktaProvider represents Okta configuration
type OktaProvider struct {
	CommonProvider `yaml:",inline"`
	OrgName        string  `yaml:"org_name,omitempty"`
	BaseURL        *string `yaml:"base_url,omitempty"`
}

// BlessProvider represents Bless ssh provider configuration
type BlessProvider struct {
	CommonProvider    `yaml:",inline"`
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `yaml:"aws_profile,omitempty"`
	AWSRegion         string   `yaml:"aws_region,omitempty"`
	RoleArn           *string  `yaml:"role_arn,omitempty"`
}

type HerokuProvider struct {
	CommonProvider `yaml:",inline"`
}

type DatadogProvider struct {
	CommonProvider `yaml:",inline"`
}

type SentryProvider struct {
	CommonProvider `yaml:",inline"`
	BaseURL        *string `yaml:"base_url,omitempty"`
}

type TfeProvider struct {
	CommonProvider `yaml:",inline"`
	Hostname       *string `yaml:"hostname,omitempty"`
}

type KubernetesProvider struct {
	CommonProvider       `yaml:",inline"`
	ClusterComponentName *string `yaml:"cluster_component_name,omitempty"`
}

type HelmProvider struct {
	CommonProvider `yaml:",inline"`
}

type KubectlProvider struct {
	CommonProvider `yaml:",inline"`
}

type GrafanaProvider struct {
	CommonProvider `yaml:",inline"`
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

	EKS             *v2.EKSConfig     `yaml:"eks,omitempty"`
	Kind            *v2.ComponentKind `yaml:"kind,omitempty"`
	ModuleSource    *string           `yaml:"module_source"`
	ModuleName      *string           `yaml:"module_name"`
	Variables       []string          `yaml:"variables"`
	ProviderAliases map[string]string `yaml:"provider_aliases"`
	Global          *Component        `yaml:"global"`
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
	if c.TFE.ProviderAliases != nil {
		tfeConfig.ProviderAliases = *c.TFE.ProviderAliases
	}
	if c.TFE.Variables != nil {
		tfeConfig.Variables = *c.TFE.Variables
	}

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
			if componentConf.ProviderAliases != nil {
				componentPlan.ProviderAliases = *componentConf.ProviderAliases
			}
			if componentConf.Variables != nil {
				componentPlan.Variables = *componentConf.Variables
			}
			componentPlan.PathToRepoRoot = "../../../../"

			componentPlan.Global = &p.Global

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
	if auth0Config != nil && (auth0Config.Enabled == nil || (auth0Config.Enabled != nil && *auth0Config.Enabled)) {
		customProvider := false
		if auth0Config.CustomProvider != nil {
			customProvider = *auth0Config.CustomProvider
		}

		source := "alexkappa/auth0"
		if auth0Config.Source != nil {
			source = *auth0Config.Source
		}

		version := "0.42.0"
		if auth0Config.Version != nil {
			version = *auth0Config.Version
		}
		auth0Plan = &Auth0Provider{
			CommonProvider: CommonProvider{
				Enabled:        auth0Config.Enabled == nil || (auth0Config.Enabled != nil && *auth0Config.Enabled),
				CustomProvider: customProvider,
				Version:        version,
			},
			Domain: *auth0Config.Domain,
		}

		providerVersions["auth0"] = ProviderVersion{
			Source:  source,
			Version: &version,
		}
	}

	var assertPlan *AssertProvider
	assertConfig := v2.ResolveAssertProvider(commons...)
	if assertConfig != nil && (assertConfig.Enabled == nil || (assertConfig.Enabled != nil && *assertConfig.Enabled)) {
		customProvider := false
		if assertConfig.CustomProvider != nil {
			customProvider = *assertConfig.CustomProvider
		}
		version := "0.0.1"
		if assertConfig.Version != nil {
			version = *assertConfig.Version
		}
		assertPlan = &AssertProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        assertConfig.Enabled == nil || (assertConfig.Enabled != nil && *assertConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["assert"] = ProviderVersion{
			Source:  "bwoznicki/assert",
			Version: &version,
		}
	}

	var githubPlan *GithubProvider
	githubConfig := v2.ResolveGithubProvider(commons...)
	if githubConfig != nil && (githubConfig.Enabled == nil || (githubConfig.Enabled != nil && *githubConfig.Enabled)) {
		customProvider := false
		if githubConfig.CustomProvider != nil {
			customProvider = *githubConfig.CustomProvider
		}

		version := "5.16.0"
		if githubConfig.Version != nil {
			version = *githubConfig.Version
		}
		githubPlan = &GithubProvider{
			Organization: *githubConfig.Organization,
			BaseURL:      githubConfig.BaseURL,
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        githubConfig.Enabled == nil || (githubConfig.Enabled != nil && *githubConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["github"] = ProviderVersion{
			Source:  "integrations/github",
			Version: &version,
		}
	}

	var snowflakePlan *SnowflakeProvider
	snowflakeConfig := v2.ResolveSnowflakeProvider(commons...)
	if snowflakeConfig != nil && (snowflakeConfig.Enabled == nil || (snowflakeConfig.Enabled != nil && *snowflakeConfig.Enabled)) {
		customProvider := false
		if snowflakeConfig.CustomProvider != nil {
			customProvider = *snowflakeConfig.CustomProvider
		}
		version := "0.55.1"
		if snowflakeConfig.Version != nil {
			version = *snowflakeConfig.Version
		}
		snowflakePlan = &SnowflakeProvider{
			Account: *snowflakeConfig.Account,
			Role:    *snowflakeConfig.Role,
			Region:  *snowflakeConfig.Region,
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        snowflakeConfig.Enabled == nil || (snowflakeConfig.Enabled != nil && *snowflakeConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["snowflake"] = ProviderVersion{
			Source:  "Snowflake-Labs/snowflake",
			Version: &version,
		}
	}

	var oktaPlan *OktaProvider
	oktaConfig := v2.ResolveOktaProvider(commons...)
	if oktaConfig != nil && (oktaConfig.Enabled == nil || (oktaConfig.Enabled != nil && *oktaConfig.Enabled)) {
		customProvider := false
		if oktaConfig.CustomProvider != nil {
			customProvider = *oktaConfig.CustomProvider
		}
		version := "3.40.0"
		if oktaConfig.Version != nil {
			version = *oktaConfig.Version
		}
		oktaPlan = &OktaProvider{
			OrgName: *oktaConfig.OrgName,
			BaseURL: oktaConfig.BaseURL,
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        oktaConfig.Enabled == nil || (oktaConfig.Enabled != nil && *oktaConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		var registryNamespace string

		if oktaConfig.RegistryNamespace != nil && *oktaConfig.RegistryNamespace != "" {
			registryNamespace = *oktaConfig.RegistryNamespace
		} else {
			registryNamespace = "oktadeveloper"
		}
		providerVersions["okta"] = ProviderVersion{
			Source:  fmt.Sprintf("%s/okta", registryNamespace),
			Version: &version,
		}
	}

	var blessPlan *BlessProvider
	blessConfig := v2.ResolveBlessProvider(commons...)
	if blessConfig != nil &&
		(blessConfig.AWSProfile != nil || blessConfig.RoleArn != nil) &&
		blessConfig.AWSRegion != nil &&
		(blessConfig.Enabled == nil || (blessConfig.Enabled != nil && *blessConfig.Enabled)) {

		customProvider := false
		if blessConfig.CustomProvider != nil {
			customProvider = *blessConfig.CustomProvider
		}
		version := "0.5.0"
		if blessConfig.Version != nil {
			version = *blessConfig.Version
		}
		blessPlan = &BlessProvider{
			AWSProfile:        blessConfig.AWSProfile,
			AWSRegion:         *blessConfig.AWSRegion,
			AdditionalRegions: blessConfig.AdditionalRegions,
			RoleArn:           blessConfig.RoleArn,
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        (blessConfig.Enabled == nil || (blessConfig.Enabled != nil && *blessConfig.Enabled)),
				CustomProvider: customProvider,
			},
		}

		providerVersions["bless"] = ProviderVersion{
			Source:  "chanzuckerberg/bless",
			Version: &version,
		}
	}

	var herokuPlan *HerokuProvider
	herokuConfig := v2.ResolveHerokuProvider(commons...)
	// Not a fan but if enabled is not there or if it explicityly says enabled true
	if herokuConfig != nil && (herokuConfig.Enabled == nil || (herokuConfig.Enabled != nil && *herokuConfig.Enabled)) {
		customProvider := false
		if herokuConfig.CustomProvider != nil {
			customProvider = *herokuConfig.CustomProvider
		}
		version := "5.1.10"
		if herokuConfig.Version != nil {
			version = *herokuConfig.Version
		}
		herokuPlan = &HerokuProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        herokuConfig.Enabled == nil || (herokuConfig.Enabled != nil && *herokuConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["heroku"] = ProviderVersion{
			Source:  "heroku/heroku",
			Version: &version,
		}
	}

	var datadogPlan *DatadogProvider
	datadogConfig := v2.ResolveDatadogProvider(commons...)
	if datadogConfig != nil && (datadogConfig.Enabled == nil || (datadogConfig.Enabled != nil && *datadogConfig.Enabled)) {
		customProvider := false
		if datadogConfig.CustomProvider != nil {
			customProvider = *datadogConfig.CustomProvider
		}
		version := "3.20.0"
		if datadogConfig.Version != nil {
			version = *datadogConfig.Version
		}
		datadogPlan = &DatadogProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        datadogConfig.Enabled == nil || (datadogConfig.Enabled != nil && *datadogConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["datadog"] = ProviderVersion{
			Source:  "datadog/datadog",
			Version: &version,
		}
	}

	pagerdutyConfig := v2.ResolvePagerdutyProvider(commons...)
	if pagerdutyConfig != nil && (pagerdutyConfig.Enabled == nil || (pagerdutyConfig.Enabled != nil && *pagerdutyConfig.Enabled)) {
		providerVersions["pagerduty"] = ProviderVersion{
			Source:  "pagerduty/pagerduty",
			Version: pagerdutyConfig.Version,
		}
	}

	opsGenieConfig := v2.ResolveOpsGenieProvider(commons...)
	if opsGenieConfig != nil && (opsGenieConfig.Enabled == nil || (opsGenieConfig.Enabled != nil && *opsGenieConfig.Enabled)) {
		providerVersions["opsgenie"] = ProviderVersion{
			Source:  "opsgenie/opsgenie",
			Version: opsGenieConfig.Version,
		}
	}

	databricksConfig := v2.ResolveDatabricksProvider(commons...)
	if databricksConfig != nil && (databricksConfig.Enabled == nil || (databricksConfig.Enabled != nil && *databricksConfig.Enabled)) {
		providerVersions["databricks"] = ProviderVersion{
			Source:  "databricks/databricks",
			Version: databricksConfig.Version,
		}
	}

	var sentryPlan *SentryProvider
	sentryConfig := v2.ResolveSentryProvider(commons...)
	if sentryConfig != nil && (sentryConfig.Enabled == nil || (sentryConfig.Enabled != nil && *sentryConfig.Enabled)) {
		customProvider := false
		if sentryConfig.CustomProvider != nil {
			customProvider = *sentryConfig.CustomProvider
		}
		version := "0.11.2"
		if sentryConfig.Version != nil {
			version = *sentryConfig.Version
		}
		sentryPlan = &SentryProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        sentryConfig.Enabled == nil || (sentryConfig.Enabled != nil && *sentryConfig.Enabled),
				CustomProvider: customProvider,
			},
			BaseURL: sentryConfig.BaseURL,
		}

		providerVersions["sentry"] = ProviderVersion{
			Source:  "jianyuan/sentry",
			Version: &version,
		}
	}

	var tfePlan *TfeProvider

	tfeConfig := v2.ResolveTfeProvider(commons...)
	if tfeConfig != nil && (tfeConfig.Enabled == nil || (tfeConfig.Enabled != nil && *tfeConfig.Enabled)) {
		customProvider := false
		if tfeConfig.CustomProvider != nil {
			customProvider = *tfeConfig.CustomProvider
		}
		version := "0.41.0"
		if tfeConfig.Version != nil {
			version = *tfeConfig.Version
		}
		tfePlan = &TfeProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        tfeConfig.Enabled == nil || (tfeConfig.Enabled != nil && *tfeConfig.Enabled),
				CustomProvider: customProvider,
			},
			Hostname: tfeConfig.Hostname,
		}

		providerVersions["tfe"] = ProviderVersion{
			Source:  "hashicorp/tfe",
			Version: &version,
		}
	}

	var k8sPlan *KubernetesProvider

	k8sConfig := v2.ResolveKubernetesProvider(commons...)
	if k8sConfig != nil && (k8sConfig.Enabled == nil || (k8sConfig.Enabled != nil && *k8sConfig.Enabled)) {
		customProvider := false
		if k8sConfig.CustomProvider != nil {
			customProvider = *k8sConfig.CustomProvider
		}
		version := "2.20.0"
		if k8sConfig.Version != nil {
			version = *k8sConfig.Version
		}
		k8sPlan = &KubernetesProvider{
			ClusterComponentName: k8sConfig.ClusterComponentName,
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        k8sConfig.Enabled == nil || (k8sConfig.Enabled != nil && *k8sConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["kubernetes"] = ProviderVersion{
			Source:  "hashicorp/kubernetes",
			Version: &version,
		}
	}

	var helmPlan *HelmProvider

	helmConfig := v2.ResolveHelmProvider(commons...)
	if helmConfig != nil && (helmConfig.Enabled == nil || (helmConfig.Enabled != nil && *helmConfig.Enabled)) {
		customProvider := false
		if helmConfig.CustomProvider != nil {
			customProvider = *helmConfig.CustomProvider
		}
		version := "2.9.0"
		if helmConfig.Version != nil {
			version = *helmConfig.Version
		}
		helmPlan = &HelmProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        helmConfig.Enabled == nil || (helmConfig.Enabled != nil && *helmConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["helm"] = ProviderVersion{
			Source:  "hashicorp/helm",
			Version: &version,
		}
	}

	var kubectlPlan *KubectlProvider

	kubcectlConfig := v2.ResolveKubectlProvider(commons...)
	if kubcectlConfig != nil && (kubcectlConfig.Enabled == nil || (kubcectlConfig.Enabled != nil && *kubcectlConfig.Enabled)) {
		customProvider := false
		if kubcectlConfig.CustomProvider != nil {
			customProvider = *kubcectlConfig.CustomProvider
		}
		version := "1.14.0"
		if kubcectlConfig.Version != nil {
			version = *kubcectlConfig.Version
		}
		kubectlPlan = &KubectlProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        kubcectlConfig.Enabled == nil || (kubcectlConfig.Enabled != nil && *kubcectlConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["kubectl"] = ProviderVersion{
			Source:  "gavinbunney/kubectl",
			Version: &version,
		}
	}

	var grafanaPlan *GrafanaProvider

	grafanaConfig := v2.ResolveGrafanaProvider(commons...)
	if grafanaConfig != nil && (grafanaConfig.Enabled == nil || (grafanaConfig.Enabled != nil && *grafanaConfig.Enabled)) {
		customProvider := false
		if grafanaConfig.CustomProvider != nil {
			customProvider = *grafanaConfig.CustomProvider
		}
		version := "1.40.1"
		if grafanaConfig.Version != nil {
			version = *grafanaConfig.Version
		}
		grafanaPlan = &GrafanaProvider{
			CommonProvider: CommonProvider{
				Version:        version,
				Enabled:        grafanaConfig.Enabled == nil || (grafanaConfig.Enabled != nil && *grafanaConfig.Enabled),
				CustomProvider: customProvider,
			},
		}

		providerVersions["grafana"] = ProviderVersion{
			Source:  "grafana/grafana",
			Version: &version,
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
		githubActionsPlan.AWSRoleName = *githubActionsConfig.AWSIAMRoleName
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

	extraTemplates := map[string]ExtraTemplate{}
	for k, v := range v2.ResolveExtraTemplates(commons...) {
		resolvedTempl := ExtraTemplate{}

		if v.Content != nil {
			resolvedTempl.Content = *v.Content
		}
		if v.Overwrite != nil {
			resolvedTempl.Overwrite = *v.Overwrite
		}
		extraTemplates[k] = ExtraTemplate{
			Overwrite: resolvedTempl.Overwrite,
			Content:   resolvedTempl.Content,
		}
	}

	return ComponentCommon{
		ExtraTemplates:           extraTemplates,
		NeedsAWSAccountsVariable: v2.ResolveAWSAccountsNeeded(commons...),
		Backend:                  backend,
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
			Helm:                   helmPlan,
			Kubectl:                kubectlPlan,
			Okta:                   oktaPlan,
			Sentry:                 sentryPlan,
			Snowflake:              snowflakePlan,
			Tfe:                    tfePlan,
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
