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

func isDifferent(config *config.Config, modules []ModuleWrapper) bool {
	return true
}
