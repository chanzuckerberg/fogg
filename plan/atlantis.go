package plan

import (
	"fmt"
	"sort"
)

type Atlantis struct {
	Enabled  bool
	Projects []AtlantisProject
}

type AtlantisProject struct {
	Name             string
	Dir              string
	TerraformVersion string
}

// buildAtlantis will walk all the components and build an atlantis plan
func (p *Plan) buildAtlantis() Atlantis {
	enabled := false
	projects := []AtlantisProject{}

	for name, acct := range p.Accounts {
		if acct.AtlantisEnabled {
			enabled = true
			p := AtlantisProject{
				Name:             fmt.Sprintf("accounts/%s", name),
				Dir:              fmt.Sprintf("terraform/accounts/%s", name),
				TerraformVersion: acct.TerraformVersion,
			}
			projects = append(projects, p)
		}
	}

	for envName, env := range p.Envs {
		for cName, c := range env.Components {
			if c.AtlantisEnabled {
				enabled = true
				p := AtlantisProject{
					Name:             fmt.Sprintf("%s/%s", envName, cName),
					Dir:              fmt.Sprintf("terraform/envs/%s/%s", envName, cName),
					TerraformVersion: c.TerraformVersion,
				}
				projects = append(projects, p)
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
	}
}
