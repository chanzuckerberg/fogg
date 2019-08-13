package examine

import (
	"strings"

	"github.com/hashicorp/terraform/config"
	"github.com/spf13/afero"
)

//Examine loads local modules and compares them to their latest version to see differences
//TODO: Comparison between local and latest
func Examine(fs afero.Fs, path string) error {
	//Collect local modules to be updated
	config, err := GetLocalModules(path)
	if err != nil && config != nil {
		return err
	}

	//Load the latest version of each module
	globalModules, err := LatestModuleVersions(fs, config)
	if err != nil {
		return err
	}

	if globalModules != nil {
	} //To silence "declared and not used" error

	return nil
}

//IsDifferent compares modules to see if there are any differences
func IsDifferent(config *config.Config, modules []ModuleWrapper) bool {
	//TODO: Empty check?
	for _, mod := range config.Modules {
		for _, modWrap := range modules {
			if !similar(mod, modWrap) {
				return true
			}
		}
	}

	return false
}

//ExamineDifferences compares the modules and finds all of the differences
func ExamineDifferences(config *config.Config, modules []ModuleWrapper) error {
	return nil
}

func similar(mod *config.Module, modWrap ModuleWrapper) bool {
	splitMod := strings.Split(mod.Source, "?")
	splitWrap := strings.Split(modWrap.moduleSource, "?")

	if strings.HasPrefix(mod.Source, modWrap.moduleSource) {
		return true
	}
	if splitMod[0] == splitWrap[0] {
		return true
	}

	return false
}

func compare(mod *config.Module, modWrap ModuleWrapper) {
}
