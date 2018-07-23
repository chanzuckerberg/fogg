package config

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestParseDefaults(t *testing.T) {
	json := `
	{
		"defaults": {
			"aws_region_backend": "us-west-2",
			"aws_region_provider": "us-west-1",
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
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "us-west-1", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)
}

func TestParse(t *testing.T) {
	f, _ := os.Open("fixtures/full.json")
	defer f.Close()
	r := bufio.NewReader(f)
	c, e := ReadConfig(r)
	assert.Nil(t, e)
	assert.NotNil(t, c.Defaults)
	assert.Equal(t, int64(1), *c.Defaults.AccountID)
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "us-west-1", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "0.1.0", c.Defaults.AWSProviderVersion)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)

	assert.NotNil(t, c.Accounts)
	assert.Len(t, c.Accounts, 2)

	assert.NotNil(t, c.Envs)
	assert.Len(t, c.Envs, 1)
	assert.Len(t, c.Envs["stage"].Components, 1)
	env := c.Envs["stage"].Components["env"]
	assert.NotNil(t, env)
	assert.NotNil(t, env.ModuleSource)
	assert.Equal(t, "github.com/foo/bar", *env.ModuleSource)

	assert.NotNil(t, c.Modules)
}

func TestJsonFailure(t *testing.T) {
	json := `foo`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)
	assert.Nil(t, c)
	assert.NotNil(t, e)
}

func TestValidation(t *testing.T) {
	json := `{}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)

	assert.NotNil(t, c)
	assert.Nil(t, e)

	e = c.Validate()
	assert.NotNil(t, e)

	_, ok := e.(*validator.InvalidValidationError)
	assert.False(t, ok)

	err, ok := e.(validator.ValidationErrors)
	assert.True(t, ok)
	assert.Len(t, err, 10)
}

func TestInitConfig(t *testing.T) {
	c := InitConfig("proj", "reg", "buck", "prof", "me@foo.example", "0.100.0", "0.99.0")
	assert.Equal(t, "prof", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "prof", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "reg", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "reg", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "0.99.0", c.Defaults.AWSProviderVersion)
	assert.Equal(t, "buck", c.Defaults.InfraBucket)
	assert.Equal(t, "me@foo.example", c.Defaults.Owner)
	assert.Equal(t, "proj", c.Defaults.Project)
	assert.Equal(t, "0.11.7", c.Defaults.TerraformVersion)
	assert.Equal(t, "0.100.0", c.Defaults.SharedInfraVersion)
}
