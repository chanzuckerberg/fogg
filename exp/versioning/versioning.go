package versioning

import (
	"fmt"

	"github.com/spf13/afero"
)

func V(fs afero.Fs) error {
	// path := "/Users/echanakira/Desktop/learning/shared-infra/terraform/envs/staging/golinks/"

	//Collect local modules to be updated
	// localModules := GetLocalModules(path)

	//Use module call source path to make http requests
	// globalModules := GetGlobalModules(localModules)

	// fmt.Printf("Modules = %v", localModules)
	fmt.Println("Hello World")
	if fs == nil {
	}
	return nil
}
