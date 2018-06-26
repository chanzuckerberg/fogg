package plan

import (
	"fmt"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
)

type account struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegion          string
	AWSRegions         []string
	InfraBucket        string
	OtherAccounts      map[string]int64
	Owner              string
	Project            string
	TerraformVersion   string
}

type module struct {
	TerraformVersion string
}

type component struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegion          string
	AWSRegions         []string
	Component          string
	Env                string
	InfraBucket        string
	OtherComponents    []string
	Owner              string
	Project            string
	TerraformVersion   string
}

type Env struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegion          string
	AWSRegions         []string
	Env                string
	InfraBucket        string
	Owner              string
	Project            string
	TerraformVersion   string

	Components map[string]component
}

type Plan struct {
	Accounts map[string]account
	Envs     map[string]Env
	Modules  map[string]module
	Version  string
}

func Eval(fs afero.Fs, configFile string) (*Plan, error) {
	c, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, err
	}
	p := &Plan{}
	// read config and validate
	// build repo plan
	v, e := util.VersionString()
	if e != nil {
		return nil, e
	}
	p.Version = v
	p.Accounts = buildAccounts(c)
	p.Envs = buildEnvs(c)
	p.Modules = buildModules(c)
	return p, nil
}

func Print(p *Plan) error {
	fmt.Printf("Version: %s\n", p.Version)
	fmt.Println("Accounts:")
	for name, account := range p.Accounts {
		fmt.Printf("\t%s:\n", name)
		if account.AccountID != nil {
			fmt.Printf("\t\taccount id: %d\n", account.AccountID)
		}
		fmt.Printf("\t\tid: %d\n", account.AccountID)

		fmt.Printf("\t\taws_profile_backend: %v\n", account.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", account.AWSProfileProvider)
		fmt.Printf("\t\taws_provider_version: %v\n", account.AWSProviderVersion)
		fmt.Printf("\t\taws_region: %v\n", account.AWSRegion)
		fmt.Printf("\t\taws_regions: %v\n", account.AWSRegions)
		fmt.Printf("\t\tinfra_bucket: %v\n", account.InfraBucket)
		fmt.Printf("\t\tname: %v\n", account.AccountName)
		fmt.Printf("\t\towner: %v\n", account.Owner)
		fmt.Printf("\t\tproject: %v\n", account.Project)
		fmt.Printf("\t\tterraform_version: %v\n", account.TerraformVersion)

		fmt.Printf("\t\tother_accounts:\n")
		for acct, id := range account.OtherAccounts {
			fmt.Printf("\t\t\t%s: %d\n", acct, id)
		}

	}

	fmt.Println("Envs:")

	for name, env := range p.Envs {
		fmt.Printf("\t%s:\n", name)
		fmt.Printf("\t\tid: %d\n", env.AccountID)

		fmt.Printf("\t\taws_profile_backend: %v\n", env.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", env.AWSProfileProvider)
		fmt.Printf("\t\taws_provider_version: %v\n", env.AWSProviderVersion)
		fmt.Printf("\t\taws_region: %v\n", env.AWSRegion)
		fmt.Printf("\t\taws_regions: %v\n", env.AWSRegions)
		fmt.Printf("\t\tenv: %v\n", env.Env)
		fmt.Printf("\t\tinfra_bucket: %v\n", env.InfraBucket)
		fmt.Printf("\t\tname: %v\n", env.AccountName)
		fmt.Printf("\t\towner: %v\n", env.Owner)
		fmt.Printf("\t\tproject: %v\n", env.Project)
		fmt.Printf("\t\tterraform_version: %v\n", env.TerraformVersion)

		fmt.Println("\t\tComponents:")

		for name, component := range env.Components {
			fmt.Printf("\t\t\t%s:\n", name)
			fmt.Printf("\t\t\t\tid: %d\n", component.AccountID)

			fmt.Printf("\t\t\t\taws_profile_backend: %v\n", component.AWSProfileBackend)
			fmt.Printf("\t\t\t\taws_profile_provider: %v\n", component.AWSProfileProvider)
			fmt.Printf("\t\t\t\taws_provider_version: %v\n", component.AWSProviderVersion)
			fmt.Printf("\t\t\t\taws_region: %v\n", component.AWSRegion)
			fmt.Printf("\t\t\t\taws_regions: %v\n", component.AWSRegions)
			fmt.Printf("\t\t\t\tinfra_bucket: %v\n", component.InfraBucket)
			fmt.Printf("\t\t\t\tname: %v\n", component.AccountName)
			fmt.Printf("\t\t\t\tother_components: %v\n", component.OtherComponents)
			fmt.Printf("\t\t\t\towner: %v\n", component.Owner)
			fmt.Printf("\t\t\t\tproject: %v\n", component.Project)
			fmt.Printf("\t\t\t\tterraform_version: %v\n", component.TerraformVersion)
		}

	}

	fmt.Println("Modules:")
	for name, module := range p.Modules {
		fmt.Printf("\t%s:\n", name)
		fmt.Printf("\t\tterraform_version: %s\n", module.TerraformVersion)
	}
	return nil
}

func buildAccounts(c *config.Config) map[string]account {
	defaults := c.Defaults
	accountPlans := make(map[string]account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := account{}

		accountPlan.AccountName = name
		accountPlan.AccountID = config.AccountID

		accountPlan.AWSRegion = resolveRequired(defaults.AWSRegion, config.AWSRegion)
		accountPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, config.AWSRegions)

		accountPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, config.AWSProfileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, config.AWSProfileProvider)
		accountPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, config.AWSProviderVersion)
		accountPlan.OtherAccounts = resolveOtherAccounts(c.Accounts, name)
		accountPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, config.TerraformVersion)
		accountPlan.InfraBucket = resolveRequired(defaults.InfraBucket, config.InfraBucket)
		accountPlan.Owner = resolveRequired(defaults.Owner, config.Owner)
		accountPlan.Project = resolveRequired(defaults.Project, config.Project)

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func buildModules(c *config.Config) map[string]module {
	modulePlans := make(map[string]module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := module{}

		modulePlan.TerraformVersion = resolveRequired(c.Defaults.TerraformVersion, conf.TerraformVersion)
		modulePlans[name] = modulePlan
	}
	return modulePlans
}

func newEnvPlan() Env {
	ep := Env{}
	ep.Components = make(map[string]component)
	return ep
}

func buildEnvs(conf *config.Config) map[string]Env {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults
	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()

		envPlan.AccountID = envConf.AccountID
		envPlan.Env = envName

		envPlan.AWSRegion = resolveRequired(defaults.AWSRegion, envConf.AWSRegion)
		envPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, envConf.AWSRegions)

		envPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, envConf.AWSProfileBackend)
		envPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, envConf.AWSProfileProvider)
		envPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, envConf.AWSProviderVersion)

		envPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, envConf.TerraformVersion)
		envPlan.InfraBucket = resolveRequired(defaults.InfraBucket, envConf.InfraBucket)
		envPlan.Owner = resolveRequired(defaults.Owner, envConf.Owner)
		envPlan.Project = resolveRequired(defaults.Project, envConf.Project)

		for componentName, componentConf := range conf.Envs[envName].Components {
			componentPlan := component{}

			componentPlan.AccountID = resolveOptionalInt(envPlan.AccountID, componentConf.AccountID)
			componentPlan.AWSRegion = resolveRequired(envPlan.AWSRegion, componentConf.AWSRegion)
			componentPlan.AWSRegions = resolveStringArray(envPlan.AWSRegions, componentConf.AWSRegions)

			componentPlan.AWSProfileBackend = resolveRequired(envPlan.AWSProfileBackend, componentConf.AWSProfileBackend)
			componentPlan.AWSProfileProvider = resolveRequired(envPlan.AWSProfileProvider, componentConf.AWSProfileProvider)
			componentPlan.AWSProviderVersion = resolveRequired(envPlan.AWSProviderVersion, componentConf.AWSProviderVersion)
			componentPlan.AccountID = resolveOptionalInt(envPlan.AccountID, componentConf.AccountID)

			componentPlan.TerraformVersion = resolveRequired(envPlan.TerraformVersion, componentConf.TerraformVersion)
			componentPlan.InfraBucket = resolveRequired(envPlan.InfraBucket, componentConf.InfraBucket)
			componentPlan.Owner = resolveRequired(envPlan.Owner, componentConf.Owner)
			componentPlan.Project = resolveRequired(envPlan.Project, componentConf.Project)

			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)

			envPlan.Components[componentName] = componentPlan
		}

		envPlans[envName] = envPlan
	}
	return envPlans
}

func otherComponentNames(components map[string]*config.Component, thisComponent string) []string {
	r := make([]string, 0)
	for componentName := range components {
		if componentName != thisComponent {
			r = append(r, componentName)
		}
	}
	return r
}

func resolveStringArray(def []string, override []string) []string {
	if override != nil {
		return override
	}
	return def
}

func resolveRequired(def string, override *string) string {
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

func resolveOtherAccounts(accounts map[string]config.Account, currentAccount string) map[string]int64 {
	other := make(map[string]int64)
	for name, account := range accounts {
		if name != currentAccount && account.AccountID != nil {
			other[name] = *account.AccountID
		}
	}
	return other
}
