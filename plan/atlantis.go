package plan

import (
	"fmt"
	"sort"
)

type Atlantis struct {
	Enabled  bool
	Projects []AtlantisProject

	// Profiles is a map of profile name -> role infro
	Profiles map[string]AWSRole
}

type AtlantisProject struct {
	Name             string `yaml:"name"`
	Dir              string `yaml:"dir"`
	TerraformVersion string `yaml:"terraform_version"`
}

type AWSRole struct {
	AccountID string `yaml:"account_id"`
	RolePath  string `yaml:"role_path"`
	RoleName  string `yaml:"role_name"`
}

// buildAtlantis will walk all the components and build an atlantis plan
func (p *Plan) buildAtlantis() Atlantis {
	// TODO This func has a lot of duplication.
	enabled := false
	projects := []AtlantisProject{}
	profiles := map[string]AWSRole{}

	if p.Global.Atlantis.Enabled {
		enabled = true
		proj := AtlantisProject{
			Name:             "global",
			Dir:              "terraform/global",
			TerraformVersion: p.Global.TerraformVersion,
		}

		projects = append(projects, proj)

		profiles[p.Global.Backend.Profile] = AWSRole{
			AccountID: *p.Global.Backend.AccountID,
			RoleName:  p.Global.Atlantis.RoleName,
			RolePath:  p.Global.Atlantis.RolePath,
		}

		if p.Global.Providers.AWS != nil {
			a := *p.Global.Providers.AWS
			profiles[a.Profile] = AWSRole{
				AccountID: *p.Global.Backend.AccountID,
				RoleName:  p.Global.Atlantis.RoleName,
				RolePath:  p.Global.Atlantis.RolePath,
			}
		}
	}

	for name, acct := range p.Accounts {
		if acct.Atlantis.Enabled {
			enabled = true
			proj := AtlantisProject{
				Name:             fmt.Sprintf("accounts/%s", name),
				Dir:              fmt.Sprintf("terraform/accounts/%s", name),
				TerraformVersion: acct.TerraformVersion,
			}
			projects = append(projects, proj)

			profiles[acct.Backend.Profile] = AWSRole{
				AccountID: *acct.Backend.AccountID,
				RoleName:  acct.Atlantis.RoleName,
				RolePath:  acct.Atlantis.RolePath,
			}

			if acct.Providers.AWS != nil {
				a := *acct.Providers.AWS
				profiles[a.Profile] = AWSRole{
					AccountID: a.AccountID.String(),
					RoleName:  acct.Atlantis.RoleName,
					RolePath:  acct.Atlantis.RolePath,
				}
			}
		}
	}

	for envName, env := range p.Envs {
		for cName, c := range env.Components {
			if c.Atlantis.Enabled {
				enabled = true
				p := AtlantisProject{
					Name:             fmt.Sprintf("%s/%s", envName, cName),
					Dir:              fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
					TerraformVersion: c.TerraformVersion,
				}
				projects = append(projects, p)

				profiles[c.Backend.Profile] = AWSRole{
					AccountID: *c.Backend.AccountID,
					RoleName:  c.Atlantis.RoleName,
					RolePath:  c.Atlantis.RolePath,
				}

				if c.Providers.AWS != nil {
					a := *c.Providers.AWS
					profiles[a.Profile] = AWSRole{
						AccountID: a.AccountID.String(),
						RoleName:  c.Atlantis.RoleName,
						RolePath:  c.Atlantis.RolePath,
					}
				}
			}
		}
	}

	// sort so that we get stable output
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	return Atlantis{
		Enabled:  enabled,
		Projects: projects,
		Profiles: profiles,
	}
}
