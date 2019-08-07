package examine

import (
	"github.com/spf13/afero"
)

func Examine(fs afero.Fs, path string) error {
	//Collect local modules to be updated
	localModules, err := GetLocalModules(fs, path)
	if err != nil{
		return err
	}
	globalModules, err := LatestModuleVersions(fs, localModules)
	if err != nil{
		return err
	}

	//TODO:(EC)Middleware to compare local and global modules
	if globalModules != nil {} //To silence "declared and not used" error

	return nil
}
