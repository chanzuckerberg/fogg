package migrations

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//Migration Defines a fogg migration and the actions that it can perform
type Migration interface {
	Description() string                      //Describes the migration taking
	Guard(afero.Fs, string) (bool, error)     //Returns true if migration is runnable, otherwise error
	Migrate(afero.Fs, string) (string, error) //Returns path to config and err
	Prompt() bool                             //Returns whether the user would like to run the migration
}

//RunMigrations cycles through a list of migrations and applies them if necessary
func RunMigrations(fs afero.Fs, configFile string, skipPrompts bool) error {
	migrations := []Migration{
		&VersionUpgradeMigration{}, //This runs first because v1 to v2 uses json.Marshal
		&JSONToYamlMigration{},
	}

	for _, migration := range migrations {
		shouldRun, err := migration.Guard(fs, configFile)
		if err != nil {
			return err
		}
		if !shouldRun{
			continue
		}

		//If the user does not want to run this migration
		if skipPrompts || !migration.Prompt(){
			logrus.Debugf("%s was skipped", migration.Description())
			continue
		}

		configFile, err = migration.Migrate(fs, configFile)
		if err != nil {
			return err
		}
		logrus.Infof("%s was successful", migration.Description())
	}
	return nil
}
