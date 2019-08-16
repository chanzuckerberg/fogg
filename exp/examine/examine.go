package examine

import (
	"github.com/spf13/afero"
)

//Examine loads local modules and compares them to their latest version to see differences
//TODO: Comparison between local and latest
func Examine(fs afero.Fs, path string) error {
	//Collect local modules to be updated
	module, err := GetLocalModules(fs, path)
	if err != nil && module != nil {
		return err
	}

	//Load the latest version of each module
	globalModules, err := LatestModuleVersions(fs, module)
	if err != nil {
		return err
	}

	if globalModules != nil {
	} //To silence "declared and not used" error

	return nil
}
