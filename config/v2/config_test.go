package v2

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	a := assert.New(t)

	b, e := ioutil.ReadFile("../testdata/empty.json")
	a.NoError(e)
	c, e := ReadConfig(b)

	a.NoError(e)

	e = c.Validate()
	a.Error(e)

	b2, e := ioutil.ReadFile("../testdata/v2_minimal_valid.json")
	a.NoError(e)
	c, e = ReadConfig(b2)

	a.NoError(e)

	e = c.Validate()
	a.NoError(e)

}
