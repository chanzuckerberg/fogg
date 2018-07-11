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
	AllAccounts        map[string]int64
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegionBackend   string
	AWSRegionProvider  string
	AWSRegions         []string
	InfraBucket        string
	Owner              string
	Project            string
	TerraformVersion   string
	SiccMode           bool
}

type Module struct {
	SiccMode         bool
	TerraformVersion string
}

type Component struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegionBackend   string
	AWSRegionProvider  string
	AWSRegions         []string
	Component          string
	Env                string
	InfraBucket        string
	OtherComponents    []string
	Owner              string
	Project            string
	SharedInfraVersion string
	TerraformVersion   string
	SiccMode           bool
}

type Env struct {
	AccountID          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSProviderVersion string
	AWSRegionBackend   string
	AWSRegionProvider  string
	AWSRegions         []string
	Env                string
	InfraBucket        string
	Owner              string
	Project            string
	TerraformVersion   string
	Type               string
	SiccMode           bool

	Components map[string]Component
}

type Plan struct {
	Accounts map[string]account
	Envs     map[string]Env
	Global   Component
	Modules  map[string]Module
	SiccMode bool
	Version  string
}

func Eval(fs afero.Fs, configFile string, siccMode bool) (*Plan, error) {
	c, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return nil, err
	}
	p := &Plan{}
	v, e := util.VersionString()
	if e != nil {
		return nil, e
	}
	p.Version = v
	p.SiccMode = siccMode
	p.Accounts = buildAccounts(c, siccMode)
	p.Envs = buildEnvs(c, siccMode)
	p.Global = buildGlobal(c, siccMode)
	p.Modules = buildModules(c, siccMode)
	return p, nil
}

func Print(p *Plan) error {
	fmt.Printf("Version: %s\n", p.Version)
	fmt.Printf("sicc mode: %t\n", p.SiccMode)
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
		fmt.Printf("\t\taws_region_backend: %v\n", account.AWSRegionBackend)
		fmt.Printf("\t\taws_region_provider: %v\n", account.AWSRegionProvider)
		fmt.Printf("\t\taws_regions: %v\n", account.AWSRegions)
		fmt.Printf("\t\tinfra_bucket: %v\n", account.InfraBucket)
		fmt.Printf("\t\tname: %v\n", account.AccountName)
		fmt.Printf("\t\towner: %v\n", account.Owner)
		fmt.Printf("\t\tproject: %v\n", account.Project)
		fmt.Printf("\t\tterraform_version: %v\n", account.TerraformVersion)

		fmt.Printf("\t\tall_accounts:\n")
		for acct, id := range account.AllAccounts {
			fmt.Printf("\t\t\t%s: %d\n", acct, id)
		}

	}

	fmt.Println("Global:")
	fmt.Printf("\tid: %d\n", p.Global.AccountID)
	fmt.Printf("\taws_profile_backend: %v\n", p.Global.AWSProfileBackend)
	fmt.Printf("\taws_profile_provider: %v\n", p.Global.AWSProfileProvider)
	fmt.Printf("\taws_provider_version: %v\n", p.Global.AWSProviderVersion)
	fmt.Printf("\taws_region_backend: %v\n", p.Global.AWSRegionBackend)
	fmt.Printf("\taws_region_provider: %v\n", p.Global.AWSRegionProvider)
	fmt.Printf("\taws_regions: %v\n", p.Global.AWSRegions)
	fmt.Printf("\tinfra_bucket: %v\n", p.Global.InfraBucket)
	fmt.Printf("\tname: %v\n", p.Global.AccountName)
	fmt.Printf("\tother_p.Globals: %v\n", p.Global.OtherComponents)
	fmt.Printf("\towner: %v\n", p.Global.Owner)
	fmt.Printf("\tproject: %v\n", p.Global.Project)
	fmt.Printf("\tterraform_version: %v\n", p.Global.TerraformVersion)

	fmt.Println("Envs:")

	for name, env := range p.Envs {
		fmt.Printf("\t%s:\n", name)
		fmt.Printf("\t\tid: %d\n", env.AccountID)

		fmt.Printf("\t\taws_profile_backend: %v\n", env.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", env.AWSProfileProvider)
		fmt.Printf("\t\taws_provider_version: %v\n", env.AWSProviderVersion)
		fmt.Printf("\t\taws_region_backend: %v\n", env.AWSRegionBackend)
		fmt.Printf("\t\taws_region_provider: %v\n", env.AWSRegionProvider)
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
			fmt.Printf("\t\t\t\taws_region_backend: %v\n", component.AWSRegionBackend)
			fmt.Printf("\t\t\t\taws_region_provider: %v\n", component.AWSRegionProvider)
			fmt.Printf("\t\t\t\taws_regions: %v\n", component.AWSRegions)
			fmt.Printf("\t\t\t\tinfra_bucket: %v\n", component.InfraBucket)
			fmt.Printf("\t\t\t\tname: %v\n", component.AccountName)
			fmt.Printf("\t\t\t\tother_components: %v\n", component.OtherComponents)
			fmt.Printf("\t\t\t\towner: %v\n", component.Owner)
			fmt.Printf("\t\t\t\tproject: %v\n", component.Project)
			fmt.Printf("\t\t\t\tshared_infra_version: %v\n", component.SharedInfraVersion)
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

func buildAccounts(c *config.Config, siccMode bool) map[string]account {
	defaults := c.Defaults
	accountPlans := make(map[string]account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := account{}

		accountPlan.AccountName = name
		accountPlan.AccountID = config.AccountID

		accountPlan.AWSRegionBackend = resolveRequired(defaults.AWSRegionBackend, config.AWSRegionBackend)
		accountPlan.AWSRegionProvider = resolveRequired(defaults.AWSRegionProvider, config.AWSRegionProvider)
		accountPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, config.AWSRegions)

		accountPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, config.AWSProfileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, config.AWSProfileProvider)
		accountPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, config.AWSProviderVersion)
		accountPlan.AllAccounts = resolveAccounts(c.Accounts)
		accountPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, config.TerraformVersion)
		accountPlan.InfraBucket = resolveRequired(defaults.InfraBucket, config.InfraBucket)
		accountPlan.Owner = resolveRequired(defaults.Owner, config.Owner)
		accountPlan.Project = resolveRequired(defaults.Project, config.Project)
		accountPlan.SiccMode = siccMode

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func buildModules(c *config.Config, siccMode bool) map[string]Module {
	modulePlans := make(map[string]Module, len(c.Modules))
	for name, conf := range c.Modules {
		modulePlan := Module{}

		modulePlan.TerraformVersion = resolveRequired(c.Defaults.TerraformVersion, conf.TerraformVersion)
		modulePlan.SiccMode = siccMode
		modulePlans[name] = modulePlan
	}
	return modulePlans
}

func newEnvPlan() Env {
	ep := Env{}
	ep.Components = make(map[string]Component)
	return ep
}

func buildGlobal(conf *config.Config, siccMode bool) Component {
	// Global just uses defaults because that's the way sicc works. We should make it directly configurable after transition.
	componentPlan := Component{}

	// TODO add accountID to defaults
	componentPlan.AWSRegionBackend = conf.Defaults.AWSRegionBackend
	componentPlan.AWSRegionProvider = conf.Defaults.AWSRegionProvider
	componentPlan.AWSRegions = conf.Defaults.AWSRegions

	componentPlan.AWSProfileBackend = conf.Defaults.AWSProfileBackend
	componentPlan.AWSProfileProvider = conf.Defaults.AWSProfileProvider
	componentPlan.AWSProviderVersion = conf.Defaults.AWSProviderVersion
	// TODO add AccountID to defaults
	// componentPlan.AccountID = conf.Defaults.AccountID

	componentPlan.TerraformVersion = conf.Defaults.TerraformVersion
	componentPlan.InfraBucket = conf.Defaults.InfraBucket
	componentPlan.Owner = conf.Defaults.Owner
	componentPlan.Project = conf.Defaults.Project

	componentPlan.Component = "global"
	componentPlan.SiccMode = siccMode
	return componentPlan
}

func buildEnvs(conf *config.Config, siccMode bool) map[string]Env {
	envPlans := make(map[string]Env, len(conf.Envs))
	defaults := conf.Defaults
	for envName, envConf := range conf.Envs {
		envPlan := newEnvPlan()

		envPlan.AccountID = envConf.AccountID
		envPlan.Env = envName

		envPlan.AWSRegionBackend = resolveRequired(defaults.AWSRegionBackend, envConf.AWSRegionBackend)
		envPlan.AWSRegionProvider = resolveRequired(defaults.AWSRegionProvider, envConf.AWSRegionProvider)
		envPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, envConf.AWSRegions)

		envPlan.AWSProfileBackend = resolveRequired(defaults.AWSProfileBackend, envConf.AWSProfileBackend)
		envPlan.AWSProfileProvider = resolveRequired(defaults.AWSProfileProvider, envConf.AWSProfileProvider)
		envPlan.AWSProviderVersion = resolveRequired(defaults.AWSProviderVersion, envConf.AWSProviderVersion)

		envPlan.TerraformVersion = resolveRequired(defaults.TerraformVersion, envConf.TerraformVersion)
		envPlan.InfraBucket = resolveRequired(defaults.InfraBucket, envConf.InfraBucket)
		envPlan.Owner = resolveRequired(defaults.Owner, envConf.Owner)
		envPlan.Project = resolveRequired(defaults.Project, envConf.Project)
		if envConf.Type != nil {
			envPlan.Type = *envConf.Type
		} else {
			envPlan.Type = "aws"
		}
		envPlan.SiccMode = siccMode

		for componentName, componentConf := range conf.Envs[envName].Components {
			componentPlan := Component{}

			componentPlan.AWSRegionBackend = resolveRequired(envPlan.AWSRegionBackend, componentConf.AWSRegionBackend)
			componentPlan.AWSRegionProvider = resolveRequired(envPlan.AWSRegionProvider, componentConf.AWSRegionProvider)
			componentPlan.AWSRegions = resolveStringArray(envPlan.AWSRegions, componentConf.AWSRegions)

			componentPlan.AWSProfileBackend = resolveRequired(envPlan.AWSProfileBackend, componentConf.AWSProfileBackend)
			componentPlan.AWSProfileProvider = resolveRequired(envPlan.AWSProfileProvider, componentConf.AWSProfileProvider)
			componentPlan.AWSProviderVersion = resolveRequired(envPlan.AWSProviderVersion, componentConf.AWSProviderVersion)
			componentPlan.AccountID = resolveOptionalInt(envPlan.AccountID, componentConf.AccountID)

			componentPlan.TerraformVersion = resolveRequired(envPlan.TerraformVersion, componentConf.TerraformVersion)
			componentPlan.InfraBucket = resolveRequired(envPlan.InfraBucket, componentConf.InfraBucket)
			componentPlan.Owner = resolveRequired(envPlan.Owner, componentConf.Owner)
			componentPlan.Project = resolveRequired(envPlan.Project, componentConf.Project)
			componentPlan.SharedInfraVersion = resolveRequired(conf.Defaults.SharedInfraVersion, componentConf.SharedInfraVersion)

			componentPlan.Env = envName
			componentPlan.Component = componentName
			componentPlan.OtherComponents = otherComponentNames(conf.Envs[envName].Components, componentName)
			componentPlan.SiccMode = siccMode

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

func resolveAccounts(accounts map[string]config.Account) map[string]int64 {
	a := make(map[string]int64)
	for name, account := range accounts {
		if account.AccountID != nil {
			a[name] = *account.AccountID
		}
	}
	return a
}
