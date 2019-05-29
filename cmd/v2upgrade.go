package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/chanzuckerberg/fogg/config"
	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	v2upgrade.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	rootCmd.AddCommand(v2upgrade)
}

var v2upgrade = &cobra.Command{
	Use:   "v2upgrade",
	Short: "Upgrades a v1 config to a v2 config",
	Long: `This command will upgrade a v1 config to a v2 config.
	Note that this is a lossy transformation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return errs.WrapInternal(err, "couldn't parse config flag")
		}

		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		bytes, version, err := config.FindConfig(fs, configFile)
		if err != nil {
			return err
		}
		switch version {
		case 1:
			c1, err := v1.ReadConfig(bytes)
			if err != nil {
				return err
			}
			c2, err := config.UpgradeConfigVersion(c1)
			if err != nil {
				return err
			}

			marshalled, err := json.MarshalIndent(c2, "", "  ")
			if err != nil {
				return errs.WrapInternal(err, "Could not serialize config to json.")
			}
			err = ioutil.WriteFile(configFile, marshalled, 0644)
			return errs.WrapInternal(err, "Could not write config to disk")
		case 2:
			logrus.Infof("config already v%d, nothing to do", version)
			return nil

		default:
			return errs.NewUserf("config version %d unrecognized", version)
		}

		config, warnings, err := readAndValidateConfig(fs, configFile)
		printWarnings(warnings)

		err = mergeConfigValidationErrors(err)
		if err != nil {
			return err
		}

		p, err := plan.Eval(config)
		if err != nil {
			return err
		}

		return plan.Print(p)
	},
}
