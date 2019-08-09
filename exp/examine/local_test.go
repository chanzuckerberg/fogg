package examine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//TODO: Move fs to versioning.go
func TestGetLocalModules(t *testing.T) {
	r := require.New(t)

	module, err := GetLocalModules("../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(module)
}
