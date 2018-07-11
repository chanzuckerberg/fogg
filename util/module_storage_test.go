package util

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadModule(t *testing.T) {
	dir, e := ioutil.TempDir("", "fogg")
	assert.Nil(t, e)

	d, e := DownloadModule(dir, "./test-module")
	assert.Nil(t, e)
	assert.NotNil(t, d)
	assert.NotEmpty(t, d)
	// TODO more asserts
}

func TestDownloadAndParseModule(t *testing.T) {
	c, e := DownloadAndParseModule("./test-module")
	assert.Nil(t, e)
	assert.NotNil(t, c)
	assert.NotNil(t, c.Variables)
	assert.NotNil(t, c.Outputs)
	assert.Len(t, c.Variables, 2)
	assert.Len(t, c.Outputs, 2)
}
