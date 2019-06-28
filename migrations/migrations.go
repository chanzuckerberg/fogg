package migrations

import (
	"github.com/spf13/afero"
)

//Migration Defines a fogg migration and the actions that it can perform
type Migration interface {
	Guard(afero.Fs, string) (bool, error)     //returns true if migration is runnable, otherwise error
	Migrate(afero.Fs, string) (string, error) //Returns path to config and err
}

//RunMigrations cycles through a list of migrations and applies them if necessary
func RunMigrations(fs afero.Fs, configFile string) error {
	migrations := []Migration{
		&VersionUpgradeMigration{},
		&JSONToYamlMigration{},
	}

	for _, migration := range migrations {
		shouldRun, err := migration.Guard(fs, configFile)
		if err != nil {
			return err
		}
		if shouldRun == false {
			continue
		}

		configFile, err = migration.Migrate(fs, configFile)
		if err != nil {
			return err
		}
	}
	return nil
}
