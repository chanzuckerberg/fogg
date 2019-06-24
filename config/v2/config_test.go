package v2

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	a := assert.New(t)

	b, e := util.TestFile("empty")
	a.NoError(e)

	fs, _, e := util.TestFs()
	a.NoError(e)
	e = afero.WriteFile(fs, "fogg.json", b, 0644)
	a.NoError(e)
	c, e := ReadConfig(b, fs, "fogg.json")
	a.NoError(e)

	w, e := c.Validate()
	a.Error(e)
	a.Len(w, 0)

	b2, e := util.TestFile("v2_minimal_valid")
	a.NoError(e)

	c, e = ReadConfig(b2, fs, "fogg.json")
	a.NoError(e)

	w, e = c.Validate()
	a.NoError(e)
	a.Len(w, 0)
}

func TestReadConfigYaml(t *testing.T) {
	a := assert.New(t)

	b2, e := util.TestFile("v2_minimal_valid_yaml")
	a.NoError(e)

	fs, _, e := util.TestFs()
	a.NoError(e)
	e = afero.WriteFile(fs, "fogg.yml", b2, 0644)
	a.NoError(e)
	c, e := ReadConfig(b2, fs, "fogg.yml")
	a.NoError(e)

	w, e := c.Validate()
	a.NoError(e)
	a.Len(w, 0)
}

func TestReadSnowflakeProvider(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("snowflake_provider")
	r.NoError(e)
	r.NotNil(b)

	fs, _, e := util.TestFs()
	r.NoError(e)
	e = afero.WriteFile(fs, "fogg.json", b, 0644)
	r.NoError(e)
	c, e := ReadConfig(b, fs, "fogg.json")
	r.NoError(e)
	r.NotNil(c)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	r.NotNil(c.Defaults.Providers)
	r.NotNil(c.Defaults.Providers.Snowflake)
	r.Equal("foo", *c.Defaults.Providers.Snowflake.Account)
	r.Equal("bar", *c.Defaults.Providers.Snowflake.Role)
	r.Equal("us-west-2", *c.Defaults.Providers.Snowflake.Region)
}

func TestReadSnowflakeProviderYaml(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("snowflake_provider_yaml")
	r.NoError(e)
	r.NotNil(b)

	fs, _, e := util.TestFs()
	r.NoError(e)
	e = afero.WriteFile(fs, "fogg.yml", b, 0644)
	r.NoError(e)
	c, e := ReadConfig(b, fs, "fogg.yml")
	r.NoError(e)
	r.NotNil(c)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	r.NotNil(c.Defaults.Providers)
	r.NotNil(c.Defaults.Providers.Snowflake)
	r.Equal("foo", *c.Defaults.Providers.Snowflake.Account)
	r.Equal("bar", *c.Defaults.Providers.Snowflake.Role)
	r.Equal("us-west-2", *c.Defaults.Providers.Snowflake.Region)
}

func TestReadBlessProvider(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("bless_provider")
	r.NoError(e)
	r.NotNil(b)

	fs, _, e := util.TestFs()
	r.NoError(e)
	e = afero.WriteFile(fs, "fogg.json", b, 0644)
	r.NoError(e)
	c, e := ReadConfig(b, fs, "fogg.json")
	r.NoError(e)
	r.NotNil(c)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	r.NotNil(c.Defaults.Providers)
	r.NotNil(c.Defaults.Providers.Bless)
	r.Equal("foofoofoo", *c.Defaults.Providers.Bless.AWSProfile)
	r.Equal("bar", *c.Defaults.Providers.Bless.AWSRegion)
	r.Equal("0.0.0", *c.Defaults.Providers.Bless.Version)
	r.Equal([]string{"a", "b"}, c.Defaults.Providers.Bless.AdditionalRegions)
}

func TestReadOktaProvider(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("okta_provider")
	r.NoError(e)
	r.NotNil(b)

	fs, _, e := util.TestFs()
	r.NoError(e)
	e = afero.WriteFile(fs, "fogg.json", b, 0644)
	r.NoError(e)
	c, e := ReadConfig(b, fs, "fogg.json")
	r.NoError(e)
	r.NotNil(c)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	r.NotNil(c.Defaults.Providers)
	r.NotNil(c.Defaults.Providers.Okta)
	r.Equal("aversion", *c.Defaults.Providers.Okta.Version)
	r.Equal("orgname", *c.Defaults.Providers.Okta.OrgName)
}

func TestReadBlessProviderYaml(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("bless_provider_yaml")
	r.NoError(e)
	r.NotNil(b)

	fs, _, e := util.TestFs()
	r.NoError(e)
	e = afero.WriteFile(fs, "fogg.yml", b, 0644)
	r.NoError(e)
	c, e := ReadConfig(b, fs, "fogg.yml")
	r.NoError(e)
	r.NotNil(c)

	w, e := c.Validate()
	r.NoError(e)
	r.Len(w, 0)

	r.NotNil(c.Defaults.Providers)
	r.NotNil(c.Defaults.Providers.Bless)
	r.Equal("foofoofoo", *c.Defaults.Providers.Bless.AWSProfile)
	r.Equal("bar", *c.Defaults.Providers.Bless.AWSRegion)
	r.Equal("0.0.0", *c.Defaults.Providers.Bless.Version)
	r.Equal([]string{"a", "b"}, c.Defaults.Providers.Bless.AdditionalRegions)
}
