//go:build !offline
// +build !offline

package util

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestDownloadModule(t *testing.T) {
	r := require.New(t)
	dir, e := ioutil.TempDir("", "fogg")
	r.Nil(e)

	pwd, e := os.Getwd()
	r.NoError(e)

	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	d, e := DownloadModule(fs, dir, "github.com/chanzuckerberg/fogg-test-module")
	r.NoError(e)
	r.NotNil(d)
	r.NotEmpty(d)
	// TODO more asserts
}

func TestDownloadAndParseModule(t *testing.T) {
	r := require.New(t)

	pwd, e := os.Getwd()
	r.NoError(e)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	downloader := MakeDownloader("github.com/chanzuckerberg/fogg-test-module")
	c, e := downloader.DownloadAndParseModule(fs)
	r.Nil(e)
	r.NotNil(c)
	r.NotNil(c.Variables)
	r.NotNil(c.Outputs)
	r.Len(c.Variables, 2)
	r.Len(c.Outputs, 2)
}
