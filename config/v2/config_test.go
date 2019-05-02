package v2

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	a := assert.New(t)

	f, _ := os.Open("testdata/empty.json")
	defer f.Close()
	r := bufio.NewReader(f)
	c, e := ReadConfig(r)

	a.NoError(e)

	e = c.Validate()
	a.Error(e)

	f2, _ := os.Open("testdata/minimal_valid.json")
	defer f2.Close()
	r = bufio.NewReader(f2)
	c, e = ReadConfig(r)

	a.NoError(e)

	e = c.Validate()
	a.NoError(e)

}
