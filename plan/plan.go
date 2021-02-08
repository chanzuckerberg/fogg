package plan

import (
	"encoding/json"
	"errors"
	"fmt"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
}

// Common represents common fields
type Common struct {
	PathToRepoRoot   string `yaml:"path_to_repo_root"`
	TerraformVersion string `yaml:"terraform_version"`
}

//ComponentCommon represents common fields for components
type ComponentCommon struct {
	Common `yaml:",inline"`

	AccountBackends   map[string]Backend      `yaml:"account_backends"`
	Accounts          map[string]*json.Number `yaml:"all_accounts"`
	Backend           Backend                 `yaml:"backend"`
	ComponentBackends map[string]Backend      `yaml:"comonent_backends"`
	Env               string                  ` yaml:"env"`
	ExtraVars         map[string]string       `yaml:"extra_vars"`
	Name              string                  `yaml:"name"`
	Owner             string                  `yaml:"owner"`
	Project           string                  `yaml:"project"`
	Providers         Providers               `yaml:"providers"`

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

type Providers struct {
	AWS                    *AWSProvider       `yaml:"aws"`
	AWSAdditionalProviders []AWSProvider      `yaml:"aws_regional_providers"`
	Github                 *GithubProvider    `yaml:"github"`
	Heroku                 *HerokuProvider    `yaml:"heroku"`
	Snowflake              *SnowflakeProvider `yaml:"snowflake"`
	Bless                  *BlessProvider     `yaml:"bless"`
	Okta                   *OktaProvider      `yaml:"okta"`
	Datadog                *DatadogProvider   `yaml:"datadog"`
	Sentry                 *SentryProvider    `yaml:"sentry"`
	Tfe                    *TfeProvider       `yaml:"tfe"`
}

//AWSProvider represents AWS provider configuration
type AWSProvider struct {
	AccountID json.Number `yaml:"account_id"`
	Alias     *string     `yaml:"alias"`
	Profile   *string     `yaml:"profile"`
	Region    string      `yaml:"region"`
	RoleArn   *string     `yaml:"role_arn"`
	Version   string      `yaml:"version"`
}

// GithubProvider represents a configuration of a github provider
type GithubProvider struct {
	Organization string  `yaml:"organization"`
	BaseURL      *string `yaml:"base_url"`
	Version      *string `yaml:"version"`
}

//SnowflakeProvider represents Snowflake DB provider configuration
type SnowflakeProvider struct {
	Account string  `yaml:"account,omitempty"`
	Role    string  `yaml:"role,omitempty"`
	Region  string  `yaml:"region,omitempty"`
	Version *string `yaml:"version,omitempty"`
}

//OktaProvider represents Okta configuration
type OktaProvider struct {
	OrgName string  `yaml:"org_name,omitempty"`
	BaseURL *string `yaml:"base_url,omitempty"`
	Version *string `yaml:"version,omitempty"`
}

//BlessProvider represents Bless ssh provider configuration
type BlessProvider struct {
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `yaml:"aws_profile,omitempty"`
	AWSRegion         string   `yaml:"aws_region,omitempty"`
	RoleArn           *string  `yaml:"role_arn,omitempty"`
	Version           *string  `yaml:"version,omitempty"`
}

type HerokuProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type DatadogProvider struct {
	Version *string `yaml:"version,omitempty"`
}

type SentryProvider struct {
	Version *string `yaml:"version,omitempty"`
	BaseURL *string `yaml:"base_url,omitempty"`
}

type TfeProvider struct {
	Enabled  bool    `yaml:"enabled,omitempty"`
	Version  *string `yaml:"version,omitempty"`
	Hostname *string `yaml:"hostname,omitempty"`
}

// BackendKind is a enum of backends we support
type BackendKind string

const (
	// BackendKindS3 is https://www.terraform.io/docs/backends/types/s3.html
	BackendKindS3     BackendKind = "s3"
	BackendKindRemote BackendKind = "remote"
)

//Backend represents a plan for configuring the terraform backend. Only one struct member can be
//non-nil at a time
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

	Account string
	Global  *Component
}

// Component is a component
type Component struct {
	ComponentCommon `yaml:",inline"`

	EKS          *v2.EKSConfig     `yaml:"eks,omitempty"`
	Kind         *v2.ComponentKind `yaml:"kind,omitempty"`
	ModuleSource *string           `yaml:"module_source"`
	ModuleName   *string           `yaml:"module_name"`
	Global       *Component        `yaml:"global"`
}

// Env is an env
type Env struct {
	Components map[string]Component `yaml:"components"`
	Env        string               `yaml:"env"`
	EKS        *v2.EKSConfig        `yaml:"eks"`
}

// TfLint containts a plan for running tflint
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
	} else {
		panic(fmt.Sprintf("Invalid backend kind of %s", componentPlan.ComponentCommon.Backend.Kind))
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

			if componentConf.Kind.GetOrDefault() == v2.ComponentKindHelmTemplate {
				componentPlan.EKS = resolveEKSConfig(envPlan.EKS, componentConf.EKS)
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

func resolveComponentCommon(commons ...v2.Common) ComponentCommon {
	var awsPlan *AWSProvider
	awsConfig := v2.ResolveAWSProvider(commons...)
	additionalProviders := []AWSProvider{}
	var roleArn *string

	if awsConfig != nil {
		if awsConfig.Role != nil {
			tmp := fmt.Sprintf("arn:aws:iam::%s:role/%s", *awsConfig.AccountID, *awsConfig.Role)
			roleArn = &tmp
		}
		awsPlan = &AWSProvider{
			AccountID: *awsConfig.AccountID,
			Profile:   awsConfig.Profile,
			Region:    *awsConfig.Region,
			RoleArn:   roleArn,
			Version:   *awsConfig.Version,
		}

		for _, r := range awsConfig.AdditionalRegions {
			// we have to take a reference here otherwise it gets overwritten by the loop
			region := r
			additionalProviders = append(additionalProviders,
				AWSProvider{
					AccountID: *awsConfig.AccountID,
					Alias:     &region,
					Profile:   awsConfig.Profile,
					Region:    region,
					RoleArn:   roleArn,
					Version:   *awsConfig.Version,
				})
		}
	}

	var githubPlan *GithubProvider
	githubConfig := v2.ResolveGithubProvider(commons...)

	if githubConfig != nil {
		githubPlan = &GithubProvider{
			Organization: *githubConfig.Organization,
			BaseURL:      githubConfig.BaseURL,
			Version:      githubConfig.Version,
		}
	}

	var snowflakePlan *SnowflakeProvider
	snowflakeConfig := v2.ResolveSnowflakeProvider(commons...)
	if snowflakeConfig != nil {
		snowflakePlan = &SnowflakeProvider{
			Account: *snowflakeConfig.Account,
			Role:    *snowflakeConfig.Role,
			Region:  *snowflakeConfig.Region,
			Version: snowflakeConfig.Version,
		}
	}

	var oktaPlan *OktaProvider
	oktaConfig := v2.ResolveOktaProvider(commons...)
	if oktaConfig != nil {
		oktaPlan = &OktaProvider{
			OrgName: *oktaConfig.OrgName,
			Version: oktaConfig.Version,
			BaseURL: oktaConfig.BaseURL,
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
			Version:           blessConfig.Version,
		}
	}

	var herokuPlan *HerokuProvider
	herokuConfig := v2.ResolveHerokuProvider(commons...)
	if herokuConfig != nil {
		herokuPlan = &HerokuProvider{
			Version: herokuConfig.Version,
		}
	}

	var datadogPlan *DatadogProvider
	datadogConfig := v2.ResolveDatadogProvider(commons...)
	if datadogConfig != nil {
		datadogPlan = &DatadogProvider{
			Version: datadogConfig.Version,
		}
	}

	var sentryPlan *SentryProvider
	sentryConfig := v2.ResolveSentryProvider(commons...)
	if sentryConfig != nil {
		sentryPlan = &SentryProvider{
			Version: sentryConfig.Version,
		}
	}

	var tfePlan *TfeProvider

	tfeConfig := v2.ResolveTfeProvider(commons...)
	if tfeConfig.Enabled != nil && *tfeConfig.Enabled {
		tfePlan = &TfeProvider{
			Enabled:  true,
			Version:  tfeConfig.Version,
			Hostname: tfeConfig.Hostname,
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

	return ComponentCommon{
		Backend: backend,
		Providers: Providers{
			AWS:                    awsPlan,
			AWSAdditionalProviders: additionalProviders,
			Github:                 githubPlan,
			Heroku:                 herokuPlan,
			Snowflake:              snowflakePlan,
			Bless:                  blessPlan,
			Okta:                   oktaPlan,
			Datadog:                datadogPlan,
			Sentry:                 sentryPlan,
			Tfe:                    tfePlan,
		},
		TfLint:          tfLintPlan,
		ExtraVars:       v2.ResolveStringMap(v2.ExtraVarsGetter, commons...),
		Owner:           v2.ResolveRequiredString(v2.OwnerGetter, commons...),
		Project:         project,
		Common:          Common{TerraformVersion: v2.ResolveRequiredString(v2.TerraformVersionGetter, commons...)},
		TravisCI:        travisPlan,
		CircleCI:        circlePlan,
		GitHubActionsCI: githubActionsPlan,
	}
}

func resolveEKSConfig(def *v2.EKSConfig, override *v2.EKSConfig) *v2.EKSConfig {
	resolved := &v2.EKSConfig{}
	if def != nil {
		resolved.ClusterName = def.ClusterName
	}
	if override != nil {
		resolved.ClusterName = override.ClusterName
	}
	return resolved
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
