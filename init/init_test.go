package init

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	a := assert.New(t)
	project := "acme"
	region := "us-west-2"
	bucket := "acme-infra"
	table := "acme"
	profile := "acme-auth"
	owner := "infra@acme.example"
	fs, _, err := util.TestFs()
	a.NoError(err)

	conf := config.InitConfig(project, region, bucket, table, profile, owner, AWSProviderVersion)
	a.NotNil(conf)
	a.Equal(config.DefaultFoggVersion, conf.Version)

	err = writeConfig(fs, conf)
	a.NoError(err)

	exists, err := afero.Exists(fs, "fogg.yml")
	a.NoError(err)
	a.True(exists)
}
