package examine

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/afero"
)

//**Local refers to any files located within your local file system**

// GetLocalModules retrieves all terraform modules within a given directory
// TODO:(EC) Define local and global modules OR rename the values
func GetLocalModules(fs afero.Fs, dir string) (*tfconfig.Module, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	module, diag := tfconfig.LoadModule(dir)
	if diag.HasErrors() {
		return nil, errs.WrapInternal(diag.Err(), "There was an issue loading the module")
	}

	return module, nil
	// return getAllModules(fs, dir)
}
