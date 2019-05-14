package v2

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	a := assert.New(t)

	b, e := util.TestFile("empty")
	a.NoError(e)
	c, e := ReadConfig(b)
	a.NoError(e)

	e = c.Validate()
	a.Error(e)

	b2, e := util.TestFile("v2_minimal_valid")
	a.NoError(e)

	c, e = ReadConfig(b2)
	a.NoError(e)

	e = c.Validate()
	a.NoError(e)

}
