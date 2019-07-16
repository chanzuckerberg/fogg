package migrations

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestBasicYamlMigrate(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.json"
	yamlMigration := JSONToYamlMigration{}
	json := []byte(`{}`)

	fs, _, err := util.TestFs()
	r.Nil(err)

	err = afero.WriteFile(fs, confPath, json, 0644)
	r.Nil(err)

	shouldRun, err := yamlMigration.Guard(fs, confPath)
	r.Nil(err)
	r.Equal(true,shouldRun)

	configFile, err := yamlMigration.Migrate(fs, confPath)
	r.Nil(err)
	r.Equal("fogg.yml", configFile)
}

func TestGuardFileType(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.txt"
	yamlMigration := JSONToYamlMigration{}
	json := []byte(`{}`)

	fs, _, err := util.TestFs()
	r.Nil(err)

	err = afero.WriteFile(fs, confPath, json, 0644)
	r.Nil(err)

	shouldRun, err := yamlMigration.Guard(fs, confPath)
	r.Error(err, " Config file .txt was not recognized")
	r.Equal(false, shouldRun)
}

func TestGuardFail(t *testing.T) {
	r := require.New(t)
	confPath := "fogg.yml"
	yamlMigration := JSONToYamlMigration{}
	json := []byte(`{}`)

	fs, _, err := util.TestFs()
	r.Nil(err)

	err = afero.WriteFile(fs, confPath, json, 0644)
	r.Nil(err)

	shouldRun, err := yamlMigration.Guard(fs, confPath)
	r.Nil(err)
	r.Equal(false, shouldRun)
}
