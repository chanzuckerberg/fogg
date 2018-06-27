package config

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	json := `
	{
		"defaults": {
			"aws_region": "us-west-2",
			"aws_profile_backend": "czi",
			"aws_profile_provider": "czi",
			"infra_s3_bucket": "the-bucket",
			"project": "test-project",
			"shared_infra_base": "../../../../",
			"terraform_version": "0.11.0"
		}
	}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)
	assert.Nil(t, e)
	assert.NotNil(t, c.Defaults)
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegion)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)
}

func TestJsonFailure(t *testing.T) {
	json := `foo`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)
	assert.Nil(t, c)
	assert.NotNil(t, e)
}
