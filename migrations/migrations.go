package migrations

import (
	"fmt"

	prompt "github.com/segmentio/go-prompt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//Migration Defines a fogg migration and the actions that it can perform
type Migration interface {
	Guard(afero.Fs, string) (bool, error)     //Returns true if migration is runnable, otherwise error
	Migrate(afero.Fs, string) (string, error) //Returns path to config and err
	Prompt() bool                             //Returns whether the user would like to run the migration
}

//RunMigrations cycles through a list of migrations and applies them if necessary
func RunMigrations(fs afero.Fs, configFile string) error {
	migrations := []Migration{
		&VersionUpgradeMigration{"Version Migration"},
		&JSONToYamlMigration{"JSON To Yaml Migration"},
	}

	skipAll := prompt.Confirm("Would you like to run all tests")

	for _, migration := range migrations {
		fmt.Printf("Running %v\n", migration)
		shouldRun, err := migration.Guard(fs, configFile)
		if err != nil {
			return err
		}
		if shouldRun == false {
			fmt.Println("Continuing")
			continue
		}

		//Ignores prompts if user chose to run all tests
		if skipAll == false {
			//Use chooses if they want to run the migration or not
			userRun := migration.Prompt()
			if userRun == false {
				fmt.Println(userRun)
				continue
			}
		}

		configFile, err = migration.Migrate(fs, configFile)
		if err != nil {
			return err
		}
		logrus.Infof("Migration was successful")
	}
	return nil
}
