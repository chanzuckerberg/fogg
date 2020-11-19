package exp

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/exp/modules"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	moduleMapCmd.Flags().String("path", "", "path to a working directory")
	moduleMapCmd.Flags().StringP("config", "c", "fogg.yml", "Use this to override the fogg config file.")

	ExpCmd.AddCommand(moduleMapCmd)
}

var moduleMapCmd = &cobra.Command{
	Use:   "module-map",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up fs
		pwd, err := os.Getwd()
		if err != nil {
			return errs.WrapUser(err, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse config flag")
		}

		openGitOrExit(fs)

		path, _ := cmd.Flags().GetString("path")
		return modules.Run(fs, configFile, path)
	},
}
