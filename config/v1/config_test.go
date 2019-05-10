package v1

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestComponentKindGetOrDefault(t *testing.T) {
	ck := ComponentKindHelmTemplate
	assert.Equal(t, string(ck.GetOrDefault()), string(ComponentKindHelmTemplate))

	var nck *ComponentKind
	assert.Equal(t, string(nck.GetOrDefault()), string(ComponentKindTerraform))

	var zck ComponentKind
	assert.Equal(t, string(zck.GetOrDefault()), string(ComponentKindTerraform))
}

func TestParseDefaults(t *testing.T) {
	json := `
	{
		"defaults": {
			"aws_region_backend": "us-west-2",
			"aws_region_provider": "us-west-1",
			"aws_profile_backend": "czi",
			"aws_profile_provider": "czi",
			"infra_s3_bucket": "the-bucket",
			"infra_dynamo_db_table": "the-table",
			"project": "test-project",
			"terraform_version": "0.11.0"
		}
	}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig([]byte(json))
	assert.NoError(t, e)
	assert.NotNil(t, c.Defaults)
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "us-west-1", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "the-table", c.Defaults.InfraDynamoTable)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)
	assert.Equal(t, true, c.Docker)
}

func TestParse(t *testing.T) {
	a := assert.New(t)
	b, e := util.TestFile("v1_full")
	a.NoError(e)
	c, e := ReadConfig(b)
	assert.Nil(t, e)
	assert.NotNil(t, c.Defaults)
	assert.Equal(t, int64(1), c.Defaults.AccountID)
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "us-west-1", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "0.1.0", c.Defaults.AWSProviderVersion)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "the-table", c.Defaults.InfraDynamoTable)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)

	assert.NotNil(t, c.Accounts)
	assert.Len(t, c.Accounts, 2)

	assert.NotNil(t, c.Envs)
	assert.Len(t, c.Envs, 1)
	assert.Len(t, c.Envs["stage"].Components, 2)
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
	b, e := ioutil.ReadAll(r)
	assert.NoError(t, e)

	c, e := ReadConfig(b)
	assert.Nil(t, c)
	assert.NotNil(t, e)
}

func TestValidation(t *testing.T) {
	json := `{}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	b, e := ioutil.ReadAll(r)
	assert.NoError(t, e)
	c, e := ReadConfig(b)

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

func TestExtraVarsValidation(t *testing.T) {
	json := `
	{
		"defaults": {
			"aws_region_backend": "us-west-2",
			"account_id": 123456789,
			"aws_region_provider": "us-west-1",
			"aws_profile_backend": "czi",
			"aws_profile_provider": "czi",
			"aws_provider_version": "czi",
			"infra_s3_bucket": "the-bucket",
			"infra_dynamo_db_table": "the-table",
			"project": "test-project",
			"owner": "test@test.com",
			"terraform_version": "0.11.0"
		}
	}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()

	c, e := ReadConfig([]byte(json))
	assert.Nil(t, e)

	e = c.Validate()
	assert.Nil(t, e)

	c.Defaults.ExtraVars = map[string]string{}
	c.Defaults.ExtraVars["env"] = "failme"
	e = c.Validate()
	assert.NotNil(t, e)
}
