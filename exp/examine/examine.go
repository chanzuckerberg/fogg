package examine

import (
	"github.com/spf13/afero"
)

// Examine loads local modules and compares them to their latest version to see differences
// TODO: Comparison between local and latest
func Examine(fs afero.Fs, path string) error {
	//Collect local modules to be updated
	module, err := GetLocalModules(fs, path)
	if err != nil {
		return err
	}
	if module == nil {
		return nil
	}
	//Load the latest version of each module
	_, err = LatestModuleVersions(fs, module)
	return err
}
