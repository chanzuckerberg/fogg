package migrations

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestUpgradeV2(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.yml"
	conf := []byte(`{"version": 1}`)
	upgradeMigration := VersionUpgradeMigration{}

	fs, _, err := util.TestFs()
	r.Nil(err)

	err = afero.WriteFile(fs, confPath, conf, 0644)
	r.Nil(err)

	shouldRun, err := upgradeMigration.Guard(fs, confPath)
	r.Equal(true, shouldRun)
	r.Nil(err)

	configFile, err := upgradeMigration.Migrate(fs, confPath)
	r.Nil(err)
	r.Equal("fogg.yml", configFile)

	_, version, err := config.FindConfig(fs, configFile)
	r.Nil(err)
	r.Equal(2, version)
}

func TestV2DoNothing(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.yml"
	conf := []byte(`{"version": 2}`)
	upgradeMigration := VersionUpgradeMigration{}

	fs, _, err := util.TestFs()
	r.Nil(err)

	err = afero.WriteFile(fs, confPath, conf, 0644)
	r.Nil(err)

	shouldRun, err := upgradeMigration.Guard(fs, confPath)
	r.Equal(false, shouldRun)
	r.Nil(err)

	configFile, err := upgradeMigration.Migrate(fs, confPath)
	r.Nil(err)
	r.Equal("fogg.yml", configFile)

	_, version, err := config.FindConfig(fs, configFile)
	r.Nil(err)
	r.Equal(2, version)
}

func TestUpgradeUnknownVersion(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.yml"
	conf := []byte(`{"version": 100}`)
	upgradeMigration := VersionUpgradeMigration{}

	fs, _, err := util.TestFs()
	r.Nil(err)

	err = afero.WriteFile(fs, confPath, conf, 0644)
	r.Nil(err)

	shouldRun, err := upgradeMigration.Guard(fs, confPath)
	r.Equal(false, shouldRun)
	r.Error(err, "Config file version was not recognized")

	configFile, err := upgradeMigration.Migrate(fs, confPath)
	r.Error(err, "config version 100 unrecognized")
	r.Equal(configFile, "fogg.yml")
}

func TestUpgradeV1(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.yml"
	upgradeMigration := VersionUpgradeMigration{}

	fs, _, err := util.TestFs()
	r.Nil(err)

	v1, err := util.TestFile("v1_full")
	r.NoError(err)

	err = afero.WriteFile(fs, confPath, v1, 0644)
	r.NoError(err)

	_, v, err := config.FindConfig(fs, confPath)
	r.NoError(err)
	r.Equal(1, v)

	shouldRun, err := upgradeMigration.Guard(fs, confPath)
	r.Nil(err)
	r.Equal(true, shouldRun)

	configFile, err := upgradeMigration.Migrate(fs, confPath)
	r.NoError(err)

	//configFile is passed to validate the Migrate returned string
	_, v, err = config.FindConfig(fs, configFile)
	r.NoError(err)
	r.Equal(2, v) // now upgraded
}
