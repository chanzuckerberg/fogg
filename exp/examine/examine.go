package examine

import (
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

//DO I COMPARE THEM RIGHT AWAY or DO I CHECK IF THERE ARE ANY DIFFERENCES AND STORE THOSE
func isDifferent(config *config.Config, modules []ModuleWrapper) bool {
	//TODO: Empty check?
	for _, mod := range config.Modules {
		for _, modWrap := range modules {
			if mod.Source == modWrap.moduleSource {
				compare(mod, modWrap)
			}
		}
	}

	return true
}

func compare(mod *config.Module, modWrap ModuleWrapper) {
	// fmt.Println(mod.Source)
	// fmt.Println(modWrap.moduleSource)
	// fmt.Printf("Module = %v\n", mod.RawConfig.Variables["var.env"])
	// fmt.Printf("ModuleWrapper = %v", modWrap.module.Variables["env"])
}
