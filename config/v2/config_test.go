package v2

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	a := assert.New(t)

	b, e := util.TestFile("empty")
	a.NoError(e)
	c, e := ReadConfig(b)
	a.NoError(e)

	w, e := c.Validate()
	a.Error(e)
	a.Len(w, 0)

	b2, e := util.TestFile("v2_minimal_valid")
	a.NoError(e)

	c, e = ReadConfig(b2)
	a.NoError(e)

	w, e = c.Validate()
	a.NoError(e)
	a.Len(w, 0)
}

func TestReadSnowflakeProvider(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("snowflake_provider")
	r.NoError(e)
	r.NotNil(b)

	c, e := ReadConfig(b)
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

	c, e := ReadConfig(b)
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
