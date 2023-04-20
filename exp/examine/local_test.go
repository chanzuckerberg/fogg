package examine

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// TODO: Move fs to versioning.go
func TestGetLocalModules(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	pwd = "../../" + pwd
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	module, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(module)
}
