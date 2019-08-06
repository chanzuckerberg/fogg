package versioning

import (
	"github.com/spf13/afero"
)

func V(fs afero.Fs) error {
	//Collect local modules to be updated
	localModules, err := GetLocalModules(fs, "/terraform/envs/staging/golinks/")
	globalModules, err := LatestModuleVersions(fs, localModules)
	if globalModules != nil || err != nil {
	}
	return nil
}
