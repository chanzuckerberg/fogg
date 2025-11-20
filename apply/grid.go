package apply

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/chanzuckerberg/fogg/config/markers"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func applyGrid(fs afero.Fs, conf *v2.Config, p *plan.Plan) error {
	guids, err := resolveGridGUIDs(fs, p)
	if err != nil {
		return err
	}

	return writeGridMarkers(fs, conf, p, guids)
}

func resolveGridGUIDs(fs afero.Fs, p *plan.Plan) (map[string]string, error) {
	guids := make(map[string]string)

	// Accounts
	for name, acct := range p.Accounts {
		if acct.Grid != nil && acct.Grid.Enabled != nil && *acct.Grid.Enabled {
			path := fmt.Sprintf("%s/accounts/%s", util.RootPath, name)
			guid, err := resolveGUID(fs, path, acct.Grid.GUID)
			if err != nil {
				return nil, err
			}
			guids[fmt.Sprintf("account:%s", name)] = guid
		}
	}

	// Envs
	for envName, env := range p.Envs {
		for compName, comp := range env.Components {
			if comp.Grid != nil && comp.Grid.Enabled != nil && *comp.Grid.Enabled {
				path := fmt.Sprintf("%s/envs/%s/%s", util.RootPath, envName, compName)
				guid, err := resolveGUID(fs, path, comp.Grid.GUID)
				if err != nil {
					return nil, err
				}
				guids[fmt.Sprintf("component:%s:%s", envName, compName)] = guid
			}
		}
	}

	return guids, nil
}

func resolveGUID(fs afero.Fs, path string, override *string) (string, error) {
	if override != nil && *override != "" {
		return *override, nil
	}

	markerPath := filepath.Join(path, ".grid-state.yaml")
	exists, err := afero.Exists(fs, markerPath)
	if err != nil {
		return "", err
	}

	if exists {
		data, err := afero.ReadFile(fs, markerPath)
		if err != nil {
			return "", err
		}
		var m markers.Marker
		if err := yaml.Unmarshal(data, &m); err != nil {
			return "", err
		}
		if m.GUID != "" {
			return m.GUID, nil
		}
	}

	return uuid.NewString(), nil
}

func writeGridMarkers(fs afero.Fs, conf *v2.Config, p *plan.Plan, guids map[string]string) error {
	// Accounts
	for name, acct := range p.Accounts {
		if acct.Grid != nil && acct.Grid.Enabled != nil && *acct.Grid.Enabled {
			path := fmt.Sprintf("%s/accounts/%s", util.RootPath, name)
			guid := guids[fmt.Sprintf("account:%s", name)]

			// Grid always uses HTTP backend, so LogicalID should be available
			if acct.ComponentCommon.Backend.HTTP == nil {
				return fmt.Errorf("account %s has Grid enabled but HTTP backend is nil", name)
			}
			logicalID := acct.ComponentCommon.Backend.HTTP.LogicalID

			var configAccount *v2.Account
			if ca, ok := conf.Accounts[name]; ok {
				configAccount = &ca
			}

			deps := resolveDependencies(acct.ComponentCommon, "", "", configAccount, nil, guids)
			if err := writeMarker(fs, path, guid, logicalID, deps); err != nil {
				return err
			}
		}
	}

	// Envs
	for envName, env := range p.Envs {
		for compName, comp := range env.Components {
			if comp.Grid != nil && comp.Grid.Enabled != nil && *comp.Grid.Enabled {
				path := fmt.Sprintf("%s/envs/%s/%s", util.RootPath, envName, compName)
				guid := guids[fmt.Sprintf("component:%s:%s", envName, compName)]

				// Grid always uses HTTP backend, so LogicalID should be available
				if comp.ComponentCommon.Backend.HTTP == nil {
					return fmt.Errorf("component %s/%s has Grid enabled but HTTP backend is nil", envName, compName)
				}
				logicalID := comp.ComponentCommon.Backend.HTTP.LogicalID

				var configComponent *v2.Component
				if e, ok := conf.Envs[envName]; ok {
					if c, ok := e.Components[compName]; ok {
						configComponent = &c
					}
				}

				deps := resolveDependencies(comp.ComponentCommon, envName, compName, nil, configComponent, guids)
				if err := writeMarker(fs, path, guid, logicalID, deps); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resolveDependencies(c plan.ComponentCommon, envName, compName string, configAccount *v2.Account, configComponent *v2.Component, guids map[string]string) []markers.Dependency {
	var deps []markers.Dependency

	// Only include dependencies if HasDependsOn is true
	if !c.HasDependsOn {
		return deps
	}

	// Get the dependency list from config
	var accountDeps, componentDeps v2.DependencyList
	if configAccount != nil && configAccount.DependsOn != nil {
		accountDeps = configAccount.DependsOn.Accounts
	}
	if configComponent != nil && configComponent.DependsOn != nil {
		accountDeps = configComponent.DependsOn.Accounts
		componentDeps = configComponent.DependsOn.Components
	}

	// DependsOn Accounts
	// AccountBackends contains only the accounts this component explicitly depends on
	// Sort the account names to ensure deterministic order
	accountNames := make([]string, 0, len(c.AccountBackends))
	for name := range c.AccountBackends {
		accountNames = append(accountNames, name)
	}
	sort.Strings(accountNames)

	for _, name := range accountNames {
		if guid, ok := guids[fmt.Sprintf("account:%s", name)]; ok {
			// Get the outputs for this specific account dependency
			outputs := getOutputsForDependency(accountDeps, name)
			// Sort outputs for deterministic order
			sort.Strings(outputs)
			for _, output := range outputs {
				deps = append(deps, markers.Dependency{
					GUID:   guid,
					Output: output,
				})
			}
			// If no outputs specified, add a dependency without output field.
			// grid-sync will default this to "default" output for backward compatibility.
			// This allows old configs with string list format to work with Grid.
			if len(outputs) == 0 {
				deps = append(deps, markers.Dependency{GUID: guid})
			}
		}
	}

	// DependsOn Components
	// ComponentBackends contains only the components this component explicitly depends on
	// Sort the component names to ensure deterministic order
	componentNames := make([]string, 0, len(c.ComponentBackends))
	for name := range c.ComponentBackends {
		componentNames = append(componentNames, name)
	}
	sort.Strings(componentNames)

	for _, name := range componentNames {
		if guid, ok := guids[fmt.Sprintf("component:%s:%s", envName, name)]; ok {
			// Get the outputs for this specific component dependency
			outputs := getOutputsForDependency(componentDeps, name)
			// Sort outputs for deterministic order
			sort.Strings(outputs)
			for _, output := range outputs {
				deps = append(deps, markers.Dependency{
					GUID:   guid,
					Output: output,
				})
			}
			// If no outputs specified, add a dependency without output field.
			// grid-sync will default this to "default" output for backward compatibility.
			// This allows old configs with string list format to work with Grid.
			if len(outputs) == 0 {
				deps = append(deps, markers.Dependency{GUID: guid})
			}
		}
	}

	return deps
}

// getOutputsForDependency retrieves the list of outputs for a specific dependency
func getOutputsForDependency(depList v2.DependencyList, depName string) []string {
	if depList == nil {
		return nil
	}
	if outputs, ok := depList[depName]; ok {
		return outputs
	}
	return nil
}

func writeMarker(fs afero.Fs, path string, guid, logicalID string, deps []markers.Dependency) error {
	marker := markers.Marker{
		GUID:         guid,
		LogicalID:    logicalID,
		Dependencies: deps,
		// Labels? We can add basic labels like "managed-by: fogg"
		Labels: map[string]string{
			"managed-by": "fogg",
		},
	}

	data, err := yaml.Marshal(marker)
	if err != nil {
		return err
	}

	markerPath := filepath.Join(path, ".grid-state.yaml")
	return afero.WriteFile(fs, markerPath, data, 0644)
}
