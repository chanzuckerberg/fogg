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
	CircleCI        CircleCIConfig        `yaml:"github_actions_ci"`
	GitHubActionsCI GitHubActionsCIConfig `yaml:"circle_ci"`
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

	Backend   Backend           `yaml:"backend"`
	ExtraVars map[string]string `yaml:"extra_vars"`
	Owner     string            `yaml:"owner"`
	Project   string            `yaml:"project"`
	Providers Providers         `yaml:"providers"`
	TfLint    TfLint            `yaml:"tf_lint"`

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

	if backend.S3 != nil && backend.S3.AccountID != nil {
		ciConfig.AWSProfiles[backend.S3.Profile] = AWSRole{
			AccountID: *backend.S3.AccountID,
			RoleName:  c.AWSRoleName,
		}
	}

	if provider != nil {
		ciConfig.AWSProfiles[provider.Profile] = AWSRole{
			AccountID: provider.AccountID.String(),
			RoleName:  c.AWSRoleName,
		}
	}
	return ciConfig
}

type Providers struct {
	AWS       *AWSProvider       `yaml:"aws"`
	Github    *GithubProvider    `yaml:"github"`
	Heroku    *HerokuProvider    `yaml:"heroku"`
	Snowflake *SnowflakeProvider `yaml:"snowflake"`
	Bless     *BlessProvider     `yaml:"bless"`
	Okta      *OktaProvider      `yaml:"okta"`
	Datadog   *DatadogProvider   `yaml:"datadog"`
}

//AWSProvider represents AWS provider configuration
type AWSProvider struct {
	AccountID         json.Number `yaml:"account_id"`
	Profile           string      `yaml:"profile"`
	Version           string      `yaml:"version"`
	Region            string      `yaml:"region"`
	AdditionalRegions []string    `yaml:"additional_regions"`
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
	OrgName string  `json:"org_name,omitempty"`
	Version *string `json:"version,omitempty"`
}

//BlessProvider represents Bless ssh provider configuration
type BlessProvider struct {
	AdditionalRegions []string `yaml:"additional_regions,omitempty"`
	AWSProfile        string   `yaml:"aws_profile,omitempty"`
	AWSRegion         string   `yaml:"aws_region,omitempty"`
	Version           *string  `yaml:"version,omitempty"`
}

type HerokuProvider struct{}

type DatadogProvider struct{}

// BackendType is a enum of backends we support
type BackendType string

const (
	// BackendTypeS3 is https://www.terraform.io/docs/backends/types/s3.html
	BackendTypeS3 BackendType = "s3"
)

//Backend represents a plan for configuring the terraform backend. Only one struct member can be
//non-nil at a time
type Backend struct {
	Type BackendType `yaml:"type"`
	S3   *S3Backend  `yaml:"s3,omitempty"`
}

// S3Backend represents aws backend configuration
type S3Backend struct {
	AccountID   *string `yaml:"account_id,omitempty"`
	AccountName string  `yaml:"account_name"`
	Profile     string  `yaml:"profile"`
	Region      string  `yaml:"region"`
	Bucket      string  `yaml:"bucket"`
	DynamoTable *string `yaml:"dynamo_table"`
}

// Module is a module
type Module struct {
	Common `yaml:",inline"`
}

// Account is an account
type Account struct {
	ComponentCommon `yaml:",inline"`

	AllAccounts map[string]*json.Number `yaml:"all_accounts"`
	AccountName string                  `yaml:"account_name"`
	Global      *Component
}

// Component is a component
type Component struct {
	ComponentCommon `yaml:",inline"`

	Accounts  map[string]Account `yaml:"accounts"` // Reference accounts for remote state
	Component string             `yaml:"component"`
	EKS       *v2.EKSConfig      `yaml:"eks,omitempty"`
	Env       string

	Kind            *v2.ComponentKind `yaml:"kind,omitempty"`
	ModuleSource    *string           `yaml:"module_source"`
	OtherComponents []string          `yaml:"other_components"`
	Global          *Component        `yaml:"global"`
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

		accountPlan.AccountName = name
		accountPlan.ComponentCommon = resolveComponentCommon(defaults.Common, acct.Common)
		accountPlan.AllAccounts = resolveAccounts(c.Accounts)
		accountPlan.PathToRepoRoot = "../../../"
		accountPlan.Global = &p.Global
		accountPlans[name] = accountPlan
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
	componentPlan.Accounts = p.Accounts
	defaults := conf.Defaults
	global := conf.Global

	componentPlan.ComponentCommon = resolveComponentCommon(defaults.Common, global.Common)

	componentPlan.Component = "global"
	componentPlan.OtherComponents = []string{}
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
			componentPlan.Accounts = p.Accounts

			componentPlan.ComponentCommon = resolveComponentCommon(defaults.Common, envConf.Common, componentConf.Common)

			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.PathToRepoRoot = "../../../../"

			componentPlan.Global = &p.Global

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans, nil
}

func resolveComponentCommon(commons ...v2.Common) ComponentCommon {
	var awsPlan *AWSProvider
	awsConfig := v2.ResolveAWSProvider(commons...)

	if awsConfig != nil {
		awsPlan = &AWSProvider{
			AccountID:         *awsConfig.AccountID,
			Profile:           *awsConfig.Profile,
			Version:           *awsConfig.Version,
			Region:            *awsConfig.Region,
			AdditionalRegions: awsConfig.AdditionalRegions,
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
		}
	}

	var blessPlan *BlessProvider
	blessConfig := v2.ResolveBlessProvider(commons...)
	if blessConfig != nil && blessConfig.AWSProfile != nil && blessConfig.AWSRegion != nil {
		blessPlan = &BlessProvider{
			AWSProfile:        *blessConfig.AWSProfile,
			AWSRegion:         *blessConfig.AWSRegion,
			AdditionalRegions: blessConfig.AdditionalRegions,
			Version:           blessConfig.Version,
		}
	}

	var herokuPlan *HerokuProvider
	herokuConfig := v2.ResolveHerokuProvider(commons...)
	if herokuConfig != nil {
		herokuPlan = &HerokuProvider{}
	}

	var datadogPlan *DatadogProvider
	datadogConfig := v2.ResolveDatadogProvider(commons...)
	if datadogConfig != nil {
		datadogPlan = &DatadogProvider{}
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

	return ComponentCommon{
		Backend: Backend{
			S3: &S3Backend{
				AccountID:   v2.ResolveOptionalString(v2.BackendAccountIdGetter, commons...),
				Region:      v2.ResolveRequiredString(v2.BackendRegionGetter, commons...),
				Profile:     v2.ResolveRequiredString(v2.BackendProfileGetter, commons...),
				Bucket:      v2.ResolveRequiredString(v2.BackendBucketGetter, commons...),
				DynamoTable: v2.ResolveOptionalString(v2.BackendDynamoTableGetter, commons...),
			},
		},
		Providers: Providers{
			AWS:       awsPlan,
			Github:    githubPlan,
			Heroku:    herokuPlan,
			Snowflake: snowflakePlan,
			Bless:     blessPlan,
			Okta:      oktaPlan,
			Datadog:   datadogPlan,
		},
		TfLint:          tfLintPlan,
		ExtraVars:       v2.ResolveStringMap(v2.ExtraVarsGetter, commons...),
		Owner:           v2.ResolveRequiredString(v2.OwnerGetter, commons...),
		Project:         v2.ResolveRequiredString(v2.ProjectGetter, commons...),
		Common:          Common{TerraformVersion: v2.ResolveRequiredString(v2.TerraformVersionGetter, commons...)},
		TravisCI:        travisPlan,
		CircleCI:        circlePlan,
		GitHubActionsCI: githubActionsPlan,
	}
}

func otherComponentNames(components map[string]v2.Component, thisComponent string) []string {
	r := make([]string, 0)
	for componentName, componentConf := range components {
		// Only set up remote state for terraform components
		if componentConf.Kind.GetOrDefault() != v2.ComponentKindTerraform {
			continue
		}
		if componentName != thisComponent {
			r = append(r, componentName)
		}
	}
	return r
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
		var acctID *json.Number
		if account.Providers != nil && account.Providers.AWS != nil {
			acctID = account.Providers.AWS.AccountID
		}
		a[name] = acctID
	}
	return a
}
