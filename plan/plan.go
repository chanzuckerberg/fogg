package plan

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/util"
	yaml "gopkg.in/yaml.v2"
)

const (
	// The version of the chanzuckerberg/terraform docker image to use
	dockerImageVersion = "0.2.1"
)

// AWSConfiguration represents aws configuration
type AWSConfiguration struct {
	AccountID          int64    `yaml:"account_id"`
	AccountName        string   `yaml:"account_name"`
	AWSProfileBackend  string   `yaml:"aws_profile_backend"`
	AWSProfileProvider string   `yaml:"aws_profile_provider"`
	AWSProviderVersion string   `yaml:"aws_provider_version"`
	AWSRegionBackend   string   `yaml:"aws_region_backend"`
	AWSRegionProvider  string   `yaml:"aws_region_provider"`
	AWSRegions         []string `yaml:"aws_regions"`
	InfraBucket        string   `yaml:"infra_bucket"`
	InfraDynamoTable   *string  `yaml:"infra_dynamo_table"`
}

// Common represents common fields
type Common struct {
	Docker             bool   `yaml:"docker"`
	DockerImageVersion string `yaml:"docker_image_version"`
	PathToRepoRoot     string `yaml:"path_to_repo_root"`
	TerraformVersion   string `yaml:"terraform_version"`
}

// Account is an account
type Account struct {
	AWSConfiguration `yaml:",inline"`
	Common           `yaml:",inline"`

	AllAccounts map[string]int64  `yaml:"all_accounts"`
	ExtraVars   map[string]string `yaml:"extra_vars"`
	Owner       string
	Project     string
	TfLint      TfLint
}

// Module is a module
type Module struct {
	Common `yaml:",inline"`
}

// Component is a component
type Component struct {
	Accounts         map[string]Account // Reference accounts for remote state
	AWSConfiguration `yaml:",inline"`
	Common           `yaml:",inline"`
	Component        string
	EKS              *v1.EKSConfig `yaml:"eks,omitempty"`
	Env              string
	ExtraVars        map[string]string `yaml:"extra_vars"`
	Kind             *v1.ComponentKind `yaml:"kind,omitempty"`
	ModuleSource     *string           `yaml:"module_source"`
	OtherComponents  []string          `yaml:"other_components"`
	Owner            string
	Project          string
	TfLint           TfLint
}

// Env is an env
type Env struct {
	AWSConfiguration `yaml:",inline"`
	Common           `yaml:",inline"`

	Components map[string]Component
	Env        string
	EKS        *v1.EKSConfig
	ExtraVars  map[string]string `yaml:"extra_vars"`
	Owner      string
	Project    string
	TfLint     TfLint
}

// TfLint containts a plan for running tflint
type TfLint struct {
	Enabled bool
}

type AWSProfile struct {
	Name string
	ID   int64
	Role string
}

// Plan represents a set of actions to take
type Plan struct {
	Accounts map[string]Account
	Envs     map[string]Env
	Global   Component
	Modules  map[string]Module
	TravisCI TravisCI
	Version  string
}

// Eval evaluates a config
func Eval(config *v2.Config) (*Plan, error) {
	p := &Plan{}
	v, e := util.VersionString()
	if e != nil {
		return nil, errs.WrapInternal(e, "unable to parse fogg version")
	}
	p.Version = v

	var err error
	p.Accounts = p.buildAccounts(config)
	p.Envs, err = p.buildEnvs(config)
	if err != nil {
		return nil, err
	}
	p.Global = p.buildGlobal(config)
	p.Modules = p.buildModules(config)

	if config.Tools.TravisCI != nil {
		p.TravisCI = p.buildTravisCI(config, v)
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

func (p *Plan) buildAccounts(c *v2.Config) map[string]Account {
	defaults := c.Defaults

	accountPlans := make(map[string]Account, len(c.Accounts))
	for name, acct := range c.Accounts {
		accountPlan := Account{}
		accountPlan.DockerImageVersion = dockerImageVersion

		accountPlan.AccountName = name

		accountPlan.AccountID = v2.ResolveRequiredInt64(v2.AWSProviderAccountIdGetter, defaults.Common, acct.Common)
		accountPlan.AWSProfileProvider = v2.ResolveRequiredString(v2.AWSProviderProfileGetter, defaults.Common, acct.Common)
		accountPlan.AWSProviderVersion = v2.ResolveRequiredString(v2.AWSProviderVersionGetter, defaults.Common, acct.Common)
		accountPlan.AWSRegionProvider = v2.ResolveRequiredString(v2.AWSProviderRegionGetter, defaults.Common, acct.Common)

		accountPlan.AWSRegions = v2.ResolveStringArray(defaults.Providers.AWS.AdditionalRegions, acct.Providers.AWS.AdditionalRegions)

		accountPlan.AWSRegionBackend = v2.ResolveRequiredString(v2.BackendRegionGetter, defaults.Common, acct.Common)
		accountPlan.AWSProfileBackend = v2.ResolveRequiredString(v2.BackendProfileGetter, defaults.Common, acct.Common)
		accountPlan.InfraBucket = v2.ResolveRequiredString(v2.BackendBucketGetter, defaults.Common, acct.Common)
		accountPlan.InfraDynamoTable = v2.ResolveOptionalString(v2.BackendDynamoTableGetter, defaults.Common, acct.Common)

		accountPlan.AllAccounts = resolveAccounts(c.Accounts)

		accountPlan.TerraformVersion = v2.ResolveRequiredString(v2.TerraformVersionGetter, defaults.Common, acct.Common)
		accountPlan.Owner = v2.ResolveRequiredString(v2.OwnerGetter, defaults.Common, acct.Common)
		accountPlan.Project = v2.ResolveRequiredString(v2.ProjectGetter, defaults.Common, acct.Common)

		accountPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, acct.ExtraVars)
		accountPlan.TfLint = resolveTfLint(c.Tools.TfLint, nil)

		accountPlan.PathToRepoRoot = "../../../"

		accountPlan.Docker = c.Docker

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func (p *Plan) buildModules(c *v2.Config) map[string]Module {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.DockerImageVersion = dockerImageVersion
		modulePlan.PathToRepoRoot = "../../../"
		modulePlan.TerraformVersion = *v2.ResolveModuleTerraformVersion(c.Defaults, conf)
		modulePlan.Docker = c.Docker
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

	componentPlan.AccountID = v2.ResolveRequiredInt64(v2.AWSProviderAccountIdGetter, defaults.Common, global.Common)
	componentPlan.AWSProfileProvider = v2.ResolveRequiredString(v2.AWSProviderProfileGetter, defaults.Common, global.Common)
	componentPlan.AWSProviderVersion = v2.ResolveRequiredString(v2.AWSProviderVersionGetter, defaults.Common, global.Common)
	componentPlan.AWSRegionProvider = v2.ResolveRequiredString(v2.AWSProviderRegionGetter, defaults.Common, global.Common)

	componentPlan.AWSRegions = v2.ResolveOptionalStringSlice(v2.AWSProviderAdditionalRegionsGetter, defaults.Common, global.Common)

	componentPlan.AWSRegionBackend = v2.ResolveRequiredString(v2.BackendRegionGetter, defaults.Common, global.Common)
	componentPlan.AWSProfileBackend = v2.ResolveRequiredString(v2.BackendProfileGetter, defaults.Common, global.Common)
	componentPlan.InfraBucket = v2.ResolveRequiredString(v2.BackendBucketGetter, defaults.Common, global.Common)
	componentPlan.InfraDynamoTable = v2.ResolveOptionalString(v2.BackendDynamoTableGetter, defaults.Common, global.Common)

	componentPlan.TerraformVersion = v2.ResolveRequiredString(v2.TerraformVersionGetter, defaults.Common, global.Common)
	componentPlan.Owner = v2.ResolveRequiredString(v2.OwnerGetter, defaults.Common, global.Common)
	componentPlan.Project = v2.ResolveRequiredString(v2.ProjectGetter, defaults.Common, global.Common)

	componentPlan.Component = "global"
	componentPlan.DockerImageVersion = dockerImageVersion
	componentPlan.OtherComponents = []string{}
	componentPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, global.ExtraVars)
	componentPlan.PathToRepoRoot = "../../"

	componentPlan.TfLint = resolveTfLint(conf.Tools.TfLint, nil)

	componentPlan.Docker = conf.Docker

	return componentPlan
}

// buildEnvs must be build after accounts
func (p *Plan) buildEnvs(conf *v2.Config) (map[string]Env, error) {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults

	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()
		envPlan.Env = envName

		envPlan.AccountID = v2.ResolveRequiredInt64(v2.AWSProviderAccountIdGetter, defaults.Common, envConf.Common)
		envPlan.AWSProfileProvider = v2.ResolveRequiredString(v2.AWSProviderProfileGetter, defaults.Common, envConf.Common)
		envPlan.AWSProviderVersion = v2.ResolveRequiredString(v2.AWSProviderVersionGetter, defaults.Common, envConf.Common)
		envPlan.AWSRegionProvider = v2.ResolveRequiredString(v2.AWSProviderRegionGetter, defaults.Common, envConf.Common)

		envPlan.AWSRegions = v2.ResolveOptionalStringSlice(v2.AWSProviderAdditionalRegionsGetter, defaults.Common, envConf.Common)

		envPlan.AWSRegionBackend = v2.ResolveRequiredString(v2.BackendRegionGetter, defaults.Common, envConf.Common)
		envPlan.AWSProfileBackend = v2.ResolveRequiredString(v2.BackendProfileGetter, defaults.Common, envConf.Common)
		envPlan.InfraBucket = v2.ResolveRequiredString(v2.BackendBucketGetter, defaults.Common, envConf.Common)
		envPlan.InfraDynamoTable = v2.ResolveOptionalString(v2.BackendDynamoTableGetter, defaults.Common, envConf.Common)

		envPlan.TerraformVersion = v2.ResolveRequiredString(v2.TerraformVersionGetter, defaults.Common, envConf.Common)
		envPlan.Owner = v2.ResolveRequiredString(v2.OwnerGetter, defaults.Common, envConf.Common)
		envPlan.Project = v2.ResolveRequiredString(v2.ProjectGetter, defaults.Common, envConf.Common)

		envPlan.ExtraVars = resolveExtraVars(defaults.ExtraVars, envConf.ExtraVars)
		envPlan.TfLint = resolveTfLint(conf.Tools.TfLint, nil)
		envPlan.Docker = conf.Docker

		for componentName, componentConf := range conf.Envs[envName].Components {
			componentPlan := Component{
				Kind: componentConf.Kind,
			}
			// reference accounts for remote state
			if _, dupe := p.Accounts[componentName]; dupe {
				return nil, errs.WrapUser(fmt.Errorf("Component %s can't have same name as account", componentName), "Invalid component name")
			}

			if componentConf.Kind.GetOrDefault() == v1.ComponentKindHelmTemplate {
				componentPlan.EKS = resolveEKSConfig(envPlan.EKS, componentConf.EKS)
			}
			componentPlan.Accounts = p.Accounts

			componentPlan.AccountID = v2.ResolveRequiredInt64(v2.AWSProviderAccountIdGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.AWSRegionProvider = v2.ResolveRequiredString(v2.AWSProviderRegionGetter, defaults.Common, envConf.Common, componentConf.Common)

			componentPlan.AWSRegions = v2.ResolveOptionalStringSlice(v2.AWSProviderAdditionalRegionsGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.AWSProfileProvider = v2.ResolveRequiredString(v2.AWSProviderProfileGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.AWSProviderVersion = v2.ResolveRequiredString(v2.AWSProviderVersionGetter, defaults.Common, envConf.Common, componentConf.Common)

			componentPlan.AWSRegionBackend = v2.ResolveRequiredString(v2.BackendRegionGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.AWSProfileBackend = v2.ResolveRequiredString(v2.BackendProfileGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.InfraBucket = v2.ResolveRequiredString(v2.BackendBucketGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.InfraDynamoTable = v2.ResolveOptionalString(v2.BackendDynamoTableGetter, defaults.Common, envConf.Common, componentConf.Common)

			componentPlan.TerraformVersion = v2.ResolveRequiredString(v2.TerraformVersionGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.Owner = v2.ResolveRequiredString(v2.OwnerGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.Project = v2.ResolveRequiredString(v2.ProjectGetter, defaults.Common, envConf.Common, componentConf.Common)
			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.DockerImageVersion = dockerImageVersion
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			componentPlan.ModuleSource = componentConf.ModuleSource
			componentPlan.ExtraVars = resolveExtraVars(envPlan.ExtraVars, componentConf.ExtraVars)
			componentPlan.PathToRepoRoot = "../../../../"

			componentPlan.TfLint = resolveTfLintComponent(envPlan.TfLint, nil)
			componentPlan.Docker = conf.Docker

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans, nil
}

func otherComponentNames(components map[string]v2.Component, thisComponent string) []string {
	r := make([]string, 0)
	for componentName, componentConf := range components {
		// Only set up remote state for terraform components
		if componentConf.Kind.GetOrDefault() != v1.ComponentKindTerraform {
			continue
		}
		if componentName != thisComponent {
			r = append(r, componentName)
		}
	}
	return r
}

func resolveEKSConfig(def *v1.EKSConfig, override *v1.EKSConfig) *v1.EKSConfig {
	resolved := &v1.EKSConfig{}
	if def != nil {
		resolved.ClusterName = def.ClusterName
	}
	if override != nil {
		resolved.ClusterName = override.ClusterName
	}
	return resolved
}

func resolveExtraVars(def map[string]string, override map[string]string) map[string]string {
	resolved := map[string]string{}
	for k, v := range def {
		resolved[k] = v
	}
	for k, v := range override {
		resolved[k] = v
	}
	return resolved
}

func resolveRequired(def *string, override *string) string {
	if override != nil && *override != "" {
		return *override
	}
	return *def
}

func resolveRequiredInt(def int64, override *int64) int64 {
	if override != nil {
		return *override
	}
	return def
}

func resolveOptionalInt(def *int64, override *int64) *int64 {
	if override != nil {
		return override
	}
	return def
}

func resolveAccounts(accounts map[string]v2.Account) map[string]int64 {
	a := make(map[string]int64)
	for name, account := range accounts {
		if account.Providers != nil && account.Providers.AWS != nil && account.Providers.AWS.AccountID != nil {
			a[name] = *account.Providers.AWS.AccountID
		}
	}
	return a
}

func resolveTfLint(def *v1.TfLint, override *v1.TfLint) TfLint {
	enabled := false
	if def != nil && def.Enabled != nil {
		enabled = *def.Enabled
	}
	if override != nil && override.Enabled != nil {
		enabled = *override.Enabled
	}
	return TfLint{
		Enabled: enabled,
	}
}

func resolveTfLintComponent(def TfLint, override *v1.TfLint) TfLint {

	enabled := def.Enabled
	if override != nil && override.Enabled != nil {
		enabled = *override.Enabled
	}
	return TfLint{
		Enabled: enabled,
	}
}
