package util

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDownloadModule(t *testing.T) {
	a := assert.New(t)
	dir, e := ioutil.TempDir("", "fogg")
	assert.Nil(t, e)

	pwd, e := os.Getwd()
	a.NoError(e)

	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	d, e := DownloadModule(fs, dir, "github.com/chanzuckerberg/fogg-test-module")
	a.NoError(e)
	a.NotNil(d)
	a.NotEmpty(d)
	// TODO more asserts
}

func TestDownloadAndParseModule(t *testing.T) {
	a := assert.New(t)

	pwd, e := os.Getwd()
	a.NoError(e)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	c, e := DownloadAndParseModule(fs, "github.com/chanzuckerberg/fogg-test-module")
	assert.Nil(t, e)
	assert.NotNil(t, c)
	assert.NotNil(t, c.Variables)
	assert.NotNil(t, c.Outputs)
	assert.Len(t, c.Variables, 2)
	assert.Len(t, c.Outputs, 2)
}
