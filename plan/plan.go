package plan

import (
	"github.com/ryanking/fogg/config"
	"github.com/ryanking/fogg/util"
	"github.com/spf13/afero"
)

type account struct {
	AccountName        string
	AWSProfileBackend  string
	AWSProfileProvider string
	AWSRegion          string
	AWSRegions         string
	InfraBucket        string
	OtherAccounts      string
	Owner              string
	Project            string
	SharedInfraPath    string
	SharedInfraVersion string
	TerraformVersion   string
}

type plan struct {
	accounts map[string]account
}

func Plan(fs afero.Fs) (*plan, error) {
	c, _ := config.FindAndReadConfig(fs)
	util.Dump(c)
	// read config and validate
	// build repo plan
	// build .sicc version plan
	buildAccounts(c)
	// build modules plan
	// build envs plan
	// walk config and apply inheritance rules
	return nil, nil
}

func buildAccounts(c *config.Config) (map[string]*account, error) {
	defaults := c.Defaults
	accountPlans := make(map[string]*account, len(c.Accounts))
	for name, config := range c.Accounts {
		accountPlan := &account{}

		accountPlan.AccountName = name

		accountPlan.AWSRegion = resolveRequired(defaults.AWSRegion, config.AWSRegion)

		// Set profiles
		profile := resolveRequired(defaults.AWSProfile, config.AWSProfile)
		profileBackend := resolveOptional(defaults.AWSProfileBackend, config.AWSProfileBackend)
		profileProvider := resolveOptional(defaults.AWSProfileBackend, config.AWSProfileBackend)

		accountPlan.AWSProfileBackend = resolveRequired(profile, profileBackend)
		accountPlan.AWSProfileProvider = resolveRequired(profile, profileProvider)

		// resolve and sort other accounts

		// fix shared infra base
		accountPlans[name] = accountPlan

	}

	return accountPlans, nil
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

// def fix_local_shared_infra_base(path):
//     if path.startswith("../"):
//         return path.replace("../", "", 1)
//     return path
