package config

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestUpgradeV2(t *testing.T) {
	r := require.New(t)
	confPath := "config"
	conf := []byte(`{"version": 1}`)
	fs, _, err := util.TestFs()
	r.Nil(err)
	err = afero.WriteFile(fs, confPath, conf, 0644)
	r.Nil(err)
	err = Upgrade(fs, confPath)
	r.Nil(err)
}

func TestUpgradeUnknownVersion(t *testing.T) {
	r := require.New(t)
	confPath := "config"
	conf := []byte(`{"version": 100}`)
	fs, _, err := util.TestFs()
	r.Nil(err)
	err = afero.WriteFile(fs, confPath, conf, 0644)
	r.Nil(err)
	err = Upgrade(fs, confPath)
	r.Error(err, "config version 100 unrecognized")
}

func TestUpgradeV1(t *testing.T) {
	r := require.New(t)
	confPath := "config"
	fs, _, err := util.TestFs()
	r.Nil(err)

	v1, err := util.TestFile("v1_full")
	r.NoError(err)

	err = afero.WriteFile(fs, confPath, v1, 0644)
	r.NoError(err)

	_, v, err := FindConfig(fs, confPath)
	r.NoError(err)
	r.Equal(1, v)

	err = Upgrade(fs, confPath)
	r.NoError(err)

	_, v, err = FindConfig(fs, confPath)
	r.NoError(err)
	r.Equal(2, v) // now upgraded
}
