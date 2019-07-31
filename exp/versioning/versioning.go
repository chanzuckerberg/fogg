package versioning

import (
	"github.com/spf13/afero"
)

func V(fs afero.Fs) error {
	// path := "/Users/echanakira/Desktop/learning/shared-infra/terraform/envs/staging/golinks/"

	//Collect local modules to be updated
	localModules, err := GetLocalModules(fs, "/terraform/envs/staging/golinks/")
	r.NoError(err)
	r.NotNil(localModules)

	globalModules, err := LatestModuleVersions(fs, localModules)
	r.NoError(err)
	r.NotNil(globalModules)
	return nil
}
