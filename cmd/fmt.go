package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	// planCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")
	rootCmd.AddCommand(fmtCmd)
}

var fmtCmd = &cobra.Command{
	Use:           "fmt",
	Short:         "format",
	Long:          "",
	SilenceErrors: true, // If we don't silence here, cobra will print them. But we want to do that in cmd/root.go
	RunE: func(cmd *cobra.Command, args []string) error {

		fs, err := pwdFs()
		if err != nil {
			return err
		}

		// if fogg.yml exists, read and format it
		if fileExists("fogg.yml") {
			c, err := config.FindAndReadConfig(fs, "fogg.yml")
			if err != nil {
				return err
			}
			err = c.Write(fs, "fogg.yml")
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func pwdFs() (afero.Fs, error) {
	// Set up fs
	pwd, e := os.Getwd()
	if e != nil {
		return nil, errs.WrapUser(e, "can't get pwd")
	}
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)
	return fs, nil
}
