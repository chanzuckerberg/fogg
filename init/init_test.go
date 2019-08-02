package init

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
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

	conf := config.InitConfig(project, region, bucket, table, profile, owner, AWSProviderVersion, foggVersion)
	a.NotNil(conf)

	err = writeConfig(fs, conf)
	a.NoError(err)
}

func TestInitVersion(t *testing.T) {
	a := assert.New(t)
	project := "acme"
	region := "us-west-2"
	bucket := "acme-infra"
	table := "acme"
	profile := "acme-auth"
	owner := "infra@acme.example"

	conf := config.InitConfig(project, region, bucket, table, profile, owner, AWSProviderVersion, foggVersion)
	a.NotNil(conf)
	// a.Nil(conf.Version)
}
