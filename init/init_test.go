package init

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	r := require.New(t)
	project := "acme"
	region := "us-west-2"
	bucket := "acme-infra"
	table := "acme"
	profile := "acme-auth"
	owner := "infra@acme.example"
	awsAccountID := "123456789"

	fs, _, err := util.TestFs()
	r.NoError(err)

	conf := config.InitConfig(&project, &region, &bucket, &table, &profile, &owner, &awsAccountID, AWSProviderVersion)
	r.NotNil(conf)
	r.Equal(config.DefaultFoggVersion, conf.Version)

	err = conf.Write(fs, "fogg.yml")
	r.NoError(err)

	exists, err := afero.Exists(fs, "fogg.yml")
	r.NoError(err)
	r.True(exists)
}
