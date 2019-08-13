package examine

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestComparison(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	config, err := GetLocalModules("../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(config)

	globalModules, err := LatestModuleVersions(fs, config)
	r.NoError(err)
	r.NotNil(globalModules)

	diff := IsDifferent(config, globalModules)
	r.Equal(true, diff)

	err = ExamineDifferences(config, globalModules)
	r.NoError(err)
}
