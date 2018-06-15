package plan

import (
	"fmt"

	"github.com/ryanking/fogg/config"
	"github.com/ryanking/fogg/util"
	"github.com/spf13/afero"
)

type account struct {
	AccountId          *int64
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSRegion          string
	AWSRegions         []string
	InfraBucket        string
	OtherAccounts      map[string]int64
	Owner              string
	Project            string
	TerraformVersion   string
}

type plan struct {
	Accounts map[string]*account
}

func Plan(fs afero.Fs) (*plan, error) {
	c, _ := config.FindAndReadConfig(fs)
	p := &plan{}
	// read config and validate
	// build repo plan
	// build .sicc version plan
	p.Accounts = buildAccounts(c)
	// build modules plan
	// build envs plan
	// walk config and apply inheritance rules
	return p, nil
}

func Print(p *plan) error {
	fmt.Println("Accounts:")
	for name, account := range p.Accounts {
		fmt.Printf("\t%s:\n", name)
		if account.AccountId != nil {
			fmt.Printf("\t\taccount id: %d\n", account.AccountId)
		}
		fmt.Printf("\t\tregions: %v\n", account.AWSRegions)
		fmt.Printf("\t\tid: %v\n", account.AccountId)
		fmt.Printf("\t\tname: %v\n", account.AccountName)
		fmt.Printf("\t\taws_profile_backend: %v\n", account.AWSProfileBackend)
		fmt.Printf("\t\taws_profile_provider: %v\n", account.AWSProfileProvider)
		fmt.Printf("\t\taws_region: %v\n", account.AWSRegion)
		fmt.Printf("\t\taws_regions: %v\n", account.AWSRegions)
		fmt.Printf("\t\tinfra_bucket: %v\n", account.InfraBucket)
		fmt.Printf("\t\tother_accounts: %v\n", account.OtherAccounts)
		fmt.Printf("\t\towner: %v\n", account.Owner)
		fmt.Printf("\t\tproject: %v\n", account.Project)
		fmt.Printf("\t\tterraform_version: %v\n", account.TerraformVersion)

	}
	return nil
}

func buildAccounts(c *config.Config) map[string]*account {
	defaults := c.Defaults
	accountPlans := make(map[string]*account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := &account{}

		accountPlan.AccountName = name
		accountPlan.AccountId = config.AccountId

		accountPlan.AWSRegion = resolveRequired(defaults.AWSRegion, config.AWSRegion)

		accountPlan.AWSRegions = resolveStringArray(defaults.AWSRegions, config.AWSRegions)

		profile := resolveRequired(defaults.AWSProfile, config.AWSProfile)
		profileBackend := resolveOptional(defaults.AWSProfileBackend, config.AWSProfileBackend)
		profileProvider := resolveOptional(defaults.AWSProfileBackend, config.AWSProfileBackend)

		accountPlan.AWSProfileBackend = resolveRequired(profile, profileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(profile, profileProvider)

		accountPlan.OtherAccounts = resolveOtherAccounts(c.Accounts, name)

		accountPlans[name] = accountPlan
	}

	return accountPlans
}

func resolveStringArray(def *[]string, override *[]string) []string {
	util.Dump(def)
	util.Dump(override)
	if override != nil {
		return *override
	}
	if def != nil {
		return *def
	}
	return []string{}
}

func resolveRequired(def string, override *string) string {
	if override != nil {
		return *override
	}
	return def
}

func resolveOptional(def *string, override *string) *string {
	if override != nil {
		return override
	}
	return def
}

func resolveOtherAccounts(accounts map[string]config.Account, currentAccount string) map[string]int64 {
	other := make(map[string]int64)
	for name, account := range accounts {
		if name != currentAccount && account.AccountId != nil {
			other[name] = *account.AccountId
		}
	}
	return other
}
